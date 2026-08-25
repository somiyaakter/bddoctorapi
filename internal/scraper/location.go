package scraper

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/gocolly/colly/v2"
)

// Location represents a division/city-level doctor directory page.
type Location struct {
	Name string
	URL  string
}

// matches URLs like https://www.doctorbangladesh.com/doctors-dhaka/
// but NOT https://www.doctorbangladesh.com/hospitals-dhaka/
var locationURLPattern = regexp.MustCompile(`^/doctors-[a-z0-9-]+/?$`)

// DiscoverLocations visits the homepage and returns every unique
// location/division doctor-directory URL found in the page.
func DiscoverLocations(c *colly.Collector, startURL string) ([]Location, error) {
	seen := make(map[string]bool)
	var locations []Location

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		href := e.Attr("href")
		abs := e.Request.AbsoluteURL(href)
		if abs == "" {
			return
		}

		u, err := url.Parse(abs)
		if err != nil {
			return
		}
		if u.Host != allowedDomain {
			return
		}

		path := u.Path
		if !strings.HasSuffix(path, "/") {
			path += "/"
		}

		if !locationURLPattern.MatchString(path) {
			return
		}

		normalized := "https://" + allowedDomain + path

		if seen[normalized] {
			return
		}
		seen[normalized] = true

		name := strings.TrimSpace(e.Text)
		if name == "" {
			// fall back to deriving the name from the slug
			slug := strings.TrimPrefix(path, "/doctors-")
			slug = strings.TrimSuffix(slug, "/")
			name = strings.Title(strings.ReplaceAll(slug, "-", " "))
		}

		locations = append(locations, Location{
			Name: name,
			URL:  normalized,
		})
	})

	var visitErr error
	c.OnError(func(r *colly.Response, err error) {
		visitErr = fmt.Errorf("failed to fetch %s: %w", r.Request.URL, err)
	})

	if err := c.Visit(startURL); err != nil {
		return nil, err
	}
	c.Wait()

	if visitErr != nil {
		return nil, visitErr
	}

	return locations, nil
}