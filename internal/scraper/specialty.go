package scraper

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/gocolly/colly/v2"
)

// Specialty represents a specialty page scoped to one location.
type Specialty struct {
	Name         string
	LocationName string
	LocationURL  string
	URL          string
}

// excludedPrefixes are nav links that would otherwise match the
// "-{locationSlug}/" suffix pattern but are not specialty pages.
var excludedPrefixes = []string{"/doctors-", "/hospitals-"}

func specialtyURLPattern(locationSlug string) *regexp.Regexp {
	return regexp.MustCompile(`^/[a-z0-9]+(?:-[a-z0-9]+)*-` + regexp.QuoteMeta(locationSlug) + `/?$`)
}

func locationSlugFromURL(locURL string) (string, error) {
	u, err := url.Parse(locURL)
	if err != nil {
		return "", err
	}
	path := strings.Trim(u.Path, "/")
	if !strings.HasPrefix(path, "doctors-") {
		return "", fmt.Errorf("unexpected location url format: %s", locURL)
	}
	return strings.TrimPrefix(path, "doctors-"), nil
}

// DiscoverSpecialties visits a single location page and returns every
// unique specialty URL scoped to that location.
func DiscoverSpecialties(c *colly.Collector, loc Location) ([]Specialty, error) {
	slug, err := locationSlugFromURL(loc.URL)
	if err != nil {
		return nil, err
	}
	pattern := specialtyURLPattern(slug)

	seen := make(map[string]bool)
	var specialties []Specialty

	// Clone gives an isolated callback/visited set per location while
	// keeping the shared AllowedDomains/Limit config from NewCollector.
	cc := c.Clone()

	cc.OnHTML("a[href]", func(e *colly.HTMLElement) {
		href := e.Attr("href")
		abs := e.Request.AbsoluteURL(href)
		if abs == "" {
			return
		}

		u, err := url.Parse(abs)
		if err != nil || u.Host != allowedDomain {
			return
		}

		path := u.Path
		if !strings.HasSuffix(path, "/") {
			path += "/"
		}

		for _, prefix := range excludedPrefixes {
			if strings.HasPrefix(path, prefix) {
				return
			}
		}

		if !pattern.MatchString(path) {
			return
		}

		normalized := "https://" + allowedDomain + path
		if seen[normalized] {
			return
		}
		seen[normalized] = true

		name := strings.TrimSpace(e.ChildText("span.lang-en"))
		if name == "" {
			name = strings.TrimSpace(e.Text)
		}

		specialties = append(specialties, Specialty{
			Name:         name,
			LocationName: loc.Name,
			LocationURL:  loc.URL,
			URL:          normalized,
		})
	})

	var visitErr error
	cc.OnError(func(r *colly.Response, err error) {
		visitErr = fmt.Errorf("failed to fetch %s: %w", r.Request.URL, err)
	})

	if err := cc.Visit(loc.URL); err != nil {
		return nil, err
	}
	cc.Wait()

	if visitErr != nil {
		return nil, visitErr
	}

	return specialties, nil
}