package main

import (
	"context"
	"fmt"
	"log"

	"datalab_api/internal/config"
	"datalab_api/internal/database"
	"datalab_api/internal/doctor"
	"datalab_api/internal/scraper"
	"datalab_api/internal/taxonomy"
)

const homepageURL = "https://www.doctorbangladesh.com/"

func main() {
	ctx := context.Background()

	// Configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	// PostgreSQL connection
	db := database.NewPostgresDB(ctx, cfg.DatabaseURL)
	defer db.Close()

	// Repositories
	repo := doctor.NewRepository(db)
	taxonomyRepo := taxonomy.NewRepository(db)

	// Colly collector
	c := scraper.NewCollector()

	// 1. Discover locations
	locations, err := scraper.DiscoverLocations(c, homepageURL)
	if err != nil {
		log.Fatalf("location discovery failed: %v", err)
	}

	fmt.Printf(
		"Discovered %d location(s)\n\n",
		len(locations),
	)

	// 2. Save locations + discover/save specialties
	totalSpecialties := 0

	for _, loc := range locations {
		locID, err := taxonomyRepo.UpsertLocation(
			ctx,
			loc.Name,
			loc.URL,
		)

		if err != nil {
			log.Printf(
				"failed to save location %s: %v",
				loc.Name,
				err,
			)
			continue
		}

		specialties, err := scraper.DiscoverSpecialties(
			c,
			loc,
		)

		if err != nil {
			log.Printf(
				"specialty discovery failed for %s: %v",
				loc.Name,
				err,
			)
			continue
		}

		fmt.Printf(
			"%s (%d specialties)\n",
			loc.Name,
			len(specialties),
		)

		for _, s := range specialties {
			_, err := taxonomyRepo.UpsertSpecialty(
				ctx,
				locID,
				s.Name,
				s.URL,
			)

			if err != nil {
				log.Printf(
					"failed to save specialty %s: %v",
					s.Name,
					err,
				)
				continue
			}
		}

		totalSpecialties += len(specialties)
	}

	fmt.Printf(
		"\nTotal specialties discovered: %d\n\n",
		totalSpecialties,
	)

	// --------------------------------------------------
	// Load taxonomy for doctor matching
	// --------------------------------------------------

	allLocations, err := taxonomyRepo.ListLocations(ctx)
	if err != nil {
		log.Fatalf(
			"failed to load locations for matching: %v",
			err,
		)
	}

	allSpecialties, err := taxonomyRepo.ListSpecialties(ctx, nil)
	if err != nil {
		log.Fatalf(
			"failed to load specialties for matching: %v",
			err,
		)
	}

	locationURLByID := make(map[int64]string)

	for _, loc := range allLocations {
		locationURLByID[loc.ID] = loc.URL
	}

	// --------------------------------------------------
	// 3. Discover doctor profile URLs
	// --------------------------------------------------

	httpClient := scraper.NewHTTPClient()

	profileURLs, err := scraper.DiscoverDoctorProfileURLs(
		httpClient,
	)

	if err != nil {
		log.Fatalf(
			"doctor profile URL discovery failed: %v",
			err,
		)
	}

	fmt.Printf(
		"Discovered %d doctor profile URL(s) via sitemap\n\n",
		len(profileURLs),
	)

	// --------------------------------------------------
	// 4. Scrape all doctor profiles
	// --------------------------------------------------

	saved, failed := 0, 0

	err = scraper.ScrapeDoctorProfiles(
		c,
		profileURLs,
		func(d doctor.Doctor) {
			id, err := repo.Upsert(ctx, d)
			if err != nil {
				log.Printf(
					"failed to save %s: %v",
					d.ProfileURL,
					err,
				)
				failed++
				return
			}

			saved++

			// Collect all chamber addresses for location matching.
			var addresses []string

			for _, ch := range d.Chambers {
				addresses = append(
					addresses,
					ch.Address,
				)
			}

			// Match doctor with locations.
			for _, loc := range taxonomy.MatchLocations(
				allLocations,
				addresses,
			) {
				if err := taxonomyRepo.LinkDoctorLocation(
					ctx,
					id,
					loc.ID,
				); err != nil {
					log.Printf(
						"failed to link doctor %d to location %s: %v",
						id,
						loc.Name,
						err,
					)
				}
			}

			// Match doctor with specialties.
			for _, sp := range taxonomy.MatchSpecialties(
				allSpecialties,
				locationURLByID,
				d.Specialties,
			) {
				if err := taxonomyRepo.LinkDoctorSpecialty(
					ctx,
					id,
					sp.ID,
				); err != nil {
					log.Printf(
						"failed to link doctor %d to specialty %s: %v",
						id,
						sp.Name,
						err,
					)
				}
			}

			log.Printf(
				"saved+linked doctor id=%d: %s",
				id,
				d.Name,
			)
		},
	)

	if err != nil {
		log.Printf(
			"scraping finished with errors: %v",
			err,
		)
	}

	// 5. Summary
	fmt.Printf(
		"\nDone. Saved: %d, Failed: %d, Attempted: %d\n",
		saved,
		failed,
		len(profileURLs),
	)
}
