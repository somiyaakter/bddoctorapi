package scraper

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const sitemapIndexURL = "https://www.doctorbangladesh.com/sitemap.xml"

// sitemapIndex mirrors the <sitemapindex> root of a sitemap index file.
type sitemapIndex struct {
	XMLName  xml.Name       `xml:"sitemapindex"`
	Sitemaps []sitemapEntry `xml:"sitemap"`
}

type sitemapEntry struct {
	Loc string `xml:"loc"`
}

// urlSet mirrors the <urlset> root of an individual sitemap file.
type urlSet struct {
	XMLName xml.Name   `xml:"urlset"`
	URLs    []urlEntry `xml:"url"`
}

type urlEntry struct {
	Loc string `xml:"loc"`
}

// fetchXML fetches url and unmarshals its XML body into v.
func fetchXML(client *http.Client, url string, v interface{}) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", NewCollector().UserAgent)

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("fetching %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status %d fetching %s", resp.StatusCode, url)
	}

	return xml.NewDecoder(resp.Body).Decode(v)
}

// DiscoverDoctorProfileURLs fetches the sitemap index and every
// "post" sub-sitemap it references, returning every doctor profile URL.
// Doctor profiles are published as WordPress "post" type content on this
// site; "page" sitemaps (About, Contact, Privacy Policy, etc.) are skipped.
func DiscoverDoctorProfileURLs(client *http.Client) ([]string, error) {
	var index sitemapIndex
	if err := fetchXML(client, sitemapIndexURL, &index); err != nil {
		return nil, fmt.Errorf("fetching sitemap index: %w", err)
	}

	var profileURLs []string
	seen := make(map[string]bool)

	for i, entry := range index.Sitemaps {
		if !strings.Contains(entry.Loc, "post-type-post") {
			continue // skip non-doctor sitemaps (static pages)
		}

		if i > 0 {
			time.Sleep(300 * time.Millisecond) // be polite between requests
		}

		var set urlSet
		if err := fetchXML(client, entry.Loc, &set); err != nil {
			return nil, fmt.Errorf("fetching sub-sitemap %s: %w", entry.Loc, err)
		}

		for _, u := range set.URLs {
			if u.Loc == "" || seen[u.Loc] {
				continue
			}
			seen[u.Loc] = true
			profileURLs = append(profileURLs, u.Loc)
		}
	}

	return profileURLs, nil
}