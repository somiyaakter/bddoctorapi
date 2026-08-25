package scraper

import (
	"datalab_api/internal/doctor"
	"datalab_api/internal/doctor_parser"

	"fmt"
	"log"

	"github.com/gocolly/colly/v2"
)

// DoctorCallback receives each successfully parsed doctor profile.
// The caller decides what to do with it (print, persist, queue, etc.) —
// this package has no knowledge of storage.
type DoctorCallback func(d doctor.Doctor)

// ScrapeDoctorProfiles visits each given doctor profile URL, parses the
// page via parser.ParseDoctorProfile, and invokes callback once per
// successfully parsed doctor. Errors on individual URLs are logged and
// do not stop the run; a summary error is returned only if URLs failed
// to even be queued.
func ScrapeDoctorProfiles(c *colly.Collector, urls []string, callback DoctorCallback) error {
	// Clone gives this run its own visited-URL set and callback
	// registrations, isolated from whatever else NewCollector's base
	// collector might be used for elsewhere.
	cc := c.Clone()

	cc.OnRequest(func(r *colly.Request) {
		log.Printf("[doctor] visiting: %s", r.URL)
	})

	cc.OnHTML("article.entry", func(e *colly.HTMLElement) {
		d := parser.ParseDoctorProfile(e)
		log.Printf("[doctor] parsed: %s (%d chamber(s)) -> %s", d.Name, len(d.Chambers), d.ProfileURL)
		callback(d)
	})

	cc.OnError(func(r *colly.Response, err error) {
		log.Printf("[doctor] request error: %s -> %v (status %d)", r.Request.URL, err, r.StatusCode)
	})

	var queueErrs []error
	for _, url := range urls {
		// Colly's default duplicate-URL protection (AllowURLRevisit=false)
		// means a repeated URL here is silently skipped by Visit, not an
		// error worth failing the run over — just log and continue.
		if err := cc.Visit(url); err != nil {
			log.Printf("[doctor] skipped %s: %v", url, err)
			queueErrs = append(queueErrs, fmt.Errorf("%s: %w", url, err))
			continue
		}
	}

	cc.Wait()

	if len(queueErrs) > 0 {
		return fmt.Errorf("%d url(s) failed to queue: %v", len(queueErrs), queueErrs)
	}
	return nil
}
