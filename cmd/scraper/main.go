package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"datalab_api/internal/config"
	"datalab_api/internal/database"
	"datalab_api/internal/doctor"
	"datalab_api/internal/scraper"
	"datalab_api/internal/taxonomy"
)

const homepageURL = "https://www.doctorbangladesh.com/"

// runHour is the local hour (Asia/Dhaka) the scraper runs at each day —
// 3 AM, when visitor traffic to both our site and the source site is lowest.
const runHour = 3

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	db := database.NewPostgresDB(ctx, cfg.DatabaseURL)
	defer db.Close()

	loc, err := time.LoadLocation("Asia/Dhaka")
	if err != nil {
		log.Printf("failed to load Asia/Dhaka timezone, falling back to UTC: %v", err)
		loc = time.UTC
	}

	// Run once immediately on startup — so a fresh deploy doesn't sit idle
	// for up to 24 hours before the first scrape.
	log.Println("running initial scrape on startup...")
	runScrape(ctx, db)

	// Then loop forever, running once every day at runHour.
	for {
		sleepUntilNextRun(loc)
		runScrape(ctx, db)
	}
}

// sleepUntilNextRun blocks until the next occurrence of runHour in loc.
func sleepUntilNextRun(loc *time.Location) {
	now := time.Now().In(loc)
	next := time.Date(now.Year(), now.Month(), now.Day(), runHour, 0, 0, 0, loc)
	if !next.After(now) {
		next = next.Add(24 * time.Hour)
	}
	d := time.Until(next)
	log.Printf("next scrape scheduled for %s (in %s)", next.Format(time.RFC3339), d.Round(time.Minute))
	time.Sleep(d)
}

// runScrape performs one full scrape-and-persist cycle. Any failure logs
// and returns early rather than crashing the process — a bad run today
// shouldn't prevent tomorrow's scheduled run from happening.
func runScrape(ctx context.Context, db *pgxpool.Pool) {
	start := time.Now()
	log.Println("=== scrape run starting ===")
	defer func() {
		log.Printf("=== scrape run finished (took %s) ===", time.Since(start).Round(time.Second))
	}()

	repo := doctor.NewRepository(db)
	taxonomyRepo := taxonomy.NewRepository(db)
	c := scraper.NewCollector()

	// 1. Discover locations
	locations, err := scraper.DiscoverLocations(c, homepageURL)
	if err != nil {
		log.Printf("location discovery failed: %v — skipping this run", err)
		return
	}
	fmt.Printf("Discovered %d location(s)\n\n", len(locations))

	// 2. Save locations + discover/save specialties
	totalSpecialties := 0
	for _, l := range locations {
		locID, err := taxonomyRepo.UpsertLocation(ctx, l.Name, l.URL)
		if err != nil {
			log.Printf("failed to save location %s: %v", l.Name, err)
			continue
		}

		specialties, err := scraper.DiscoverSpecialties(c, l)
		if err != nil {
			log.Printf("specialty discovery failed for %s: %v", l.Name, err)
			continue
		}
		fmt.Printf("%s (%d specialties)\n", l.Name, len(specialties))

		for _, s := range specialties {
			if _, err := taxonomyRepo.UpsertSpecialty(ctx, locID, s.Name, s.URL); err != nil {
				log.Printf("failed to save specialty %s: %v", s.Name, err)
			}
		}
		totalSpecialties += len(specialties)
	}
	fmt.Printf("\nTotal specialties discovered: %d\n\n", totalSpecialties)

	// Load taxonomy for doctor matching
	allLocations, err := taxonomyRepo.ListLocations(ctx)
	if err != nil {
		log.Printf("failed to load locations for matching: %v — skipping this run", err)
		return
	}
	allSpecialties, err := taxonomyRepo.ListSpecialties(ctx, nil)
	if err != nil {
		log.Printf("failed to load specialties for matching: %v — skipping this run", err)
		return
	}
	locationURLByID := make(map[int64]string)
	for _, l := range allLocations {
		locationURLByID[l.ID] = l.URL
	}

	// 3. Discover doctor profile URLs
	httpClient := scraper.NewHTTPClient()
	profileURLs, err := scraper.DiscoverDoctorProfileURLs(httpClient)
	if err != nil {
		log.Printf("doctor profile URL discovery failed: %v — skipping this run", err)
		return
	}
	fmt.Printf("Discovered %d doctor profile URL(s) via sitemap\n\n", len(profileURLs))

	// 4. Scrape all doctor profiles
	saved, failed := 0, 0
	err = scraper.ScrapeDoctorProfiles(c, profileURLs, func(d doctor.Doctor) {
		id, err := repo.Upsert(ctx, d)
		if err != nil {
			log.Printf("failed to save %s: %v", d.ProfileURL, err)
			failed++
			return
		}
		saved++

		var addresses []string
		for _, ch := range d.Chambers {
			addresses = append(addresses, ch.Address)
		}

		for _, l := range taxonomy.MatchLocations(allLocations, addresses) {
			if err := taxonomyRepo.LinkDoctorLocation(ctx, id, l.ID); err != nil {
				log.Printf("failed to link doctor %d to location %s: %v", id, l.Name, err)
			}
		}
		for _, sp := range taxonomy.MatchSpecialties(allSpecialties, locationURLByID, d.Specialties) {
			if err := taxonomyRepo.LinkDoctorSpecialty(ctx, id, sp.ID); err != nil {
				log.Printf("failed to link doctor %d to specialty %s: %v", id, sp.Name, err)
			}
		}

		log.Printf("saved+linked doctor id=%d: %s", id, d.Name)
	})

	if err != nil {
		log.Printf("scraping finished with errors: %v", err)
	}

	fmt.Printf("\nDone. Saved: %d, Failed: %d, Attempted: %d\n", saved, failed, len(profileURLs))
}