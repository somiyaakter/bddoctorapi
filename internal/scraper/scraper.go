package scraper

import (
	"net/http"
	"time"
	"github.com/gocolly/colly/v2"
)

const (
	allowedDomain = "www.doctorbangladesh.com"
)
// NewHTTPClient returns a plain HTTP client for non-HTML fetches
// (e.g. XML sitemaps), using the same timeout convention as the collector.
func NewHTTPClient() *http.Client {
	return &http.Client{Timeout: 15 * time.Second}
}
// NewCollector returns a Colly collector configured with sane,
// respectful defaults shared across all scraping phases.



func NewCollector() *colly.Collector {
	c := colly.NewCollector(
		colly.AllowedDomains(allowedDomain),
	)

	c.SetRequestTimeout(15 * time.Second)

	
	_ = c.Limit(&colly.LimitRule{
		DomainGlob:  "*doctorbangladesh.*",
		Parallelism: 5,
		Delay:       500 * time.Millisecond,
	})

	return c
}
