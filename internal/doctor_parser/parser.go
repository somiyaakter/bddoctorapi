package parser

import (
	"datalab_api/internal/doctor"
	"html"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
)

var (
	whitespaceRegex         = regexp.MustCompile(`\s+`)
	experienceYearsRegex    = regexp.MustCompile(`(\d+)`)
	bmdcRegex               = regexp.MustCompile(`(?i)BMDC[^:]*:\s*(.+)`)
	chamberHeadingRegex     = regexp.MustCompile(`(?i)chamber.*appointment`)
	brSplitRegex            = regexp.MustCompile(`(?i)<br\s*/?>`)
	tagStripRegex           = regexp.MustCompile(`<[^>]+>`)
	addressPrefixRegex      = regexp.MustCompile(`(?i)^address\s*:\s*`)
	visitingHourPrefixRegex = regexp.MustCompile(`(?i)^visiting\s*hour\s*:\s*`)
	appointmentPrefixRegex  = regexp.MustCompile(`(?i)^appointment\s*:\s*`)
)

// cleanText collapses newlines/tabs/repeated whitespace into single spaces
// and trims the result, per the project's HTML-cleaning requirement.
func cleanText(text string) string {
	return strings.TrimSpace(whitespaceRegex.ReplaceAllString(text, " "))
}

// stripTags removes HTML tags from a fragment and decodes entities
// (e.g. "&amp;" -> "&"), used when parsing chamber <p> blocks split on <br>.
func stripTags(s string) string {
	return cleanText(html.UnescapeString(tagStripRegex.ReplaceAllString(s, "")))
}

// parseExperienceYears pulls the leading integer out of text like
// "17+ Years of Experience". Returns 0 if no number is present (the
// li can be empty, per the original Section 6 example HTML).
func parseExperienceYears(text string) int {
	match := experienceYearsRegex.FindString(text)
	if match == "" {
		return 0
	}
	years, err := strconv.Atoi(match)
	if err != nil {
		return 0
	}
	return years
}

// extractBMDCRegNo finds the BMDC registration number. On the live site
// this li has NO title attribute (unlike Degree/Specialty/Designation/
// Workplace, which all do) — so we identify it by that absence plus a
// "BMDC" text match, rather than by position. Returns "" if not present;
// BMDC is treated as optional since not every profile shows it.
func extractBMDCRegNo(doctorInfo *goquery.Selection) string {
	var regNo string
	doctorInfo.Find("li").EachWithBreak(func(_ int, li *goquery.Selection) bool {
		if _, hasTitle := li.Attr("title"); hasTitle {
			return true // a labeled field (Degree/Specialty/etc) — keep looking
		}
		text := cleanText(li.Text())
		if m := bmdcRegex.FindStringSubmatch(text); m != nil {
			regNo = cleanText(m[1])
			return false // found it, stop
		}
		return true
	})
	return regNo
}

// parseChamberParagraph extracts one chamber's data from the <p> that
// immediately follows a "Chamber [N] & Appointment" heading. The <p> is
// a flat run of inline elements separated by <br>, e.g.:
//
//	<strong><a href="...">Hospital Name</a></strong><br>
//	Address: ...<br>
//	<strong>Visiting Hour: ...</strong><br>
//	Appointment: ...<br>
//	<a class="call-now" href="tel:...">Call Now</a>
//
// so we split the raw inner HTML on <br> and classify each resulting
// line by its "Address:" / "Visiting Hour:" / "Appointment:" prefix.
func parseChamberParagraph(p *goquery.Selection) doctor.Chamber {
	var chamber doctor.Chamber

	if link := p.Find("strong a").First(); link.Length() > 0 {
		chamber.Name = cleanText(link.Text())
	}

	innerHTML, _ := p.Html()
	for _, raw := range brSplitRegex.Split(innerHTML, -1) {
		line := stripTags(raw)
		if line == "" {
			continue
		}
		switch {
		case addressPrefixRegex.MatchString(line):
			chamber.Address = addressPrefixRegex.ReplaceAllString(line, "")
		case visitingHourPrefixRegex.MatchString(line):
			chamber.VisitingHour = visitingHourPrefixRegex.ReplaceAllString(line, "")
		case appointmentPrefixRegex.MatchString(line):
			chamber.AppointmentPhone = appointmentPrefixRegex.ReplaceAllString(line, "")
		}
	}

	return chamber
}

// ParseDoctorProfile extracts a Doctor (with its Chambers) from a doctor
// profile page. e.Request.URL is used as the ProfileURL — never scraped
// from page content, per the project's rule.
func ParseDoctorProfile(e *colly.HTMLElement) doctor.Doctor {
	doctorInfo := e.DOM.Find("ul.doctor-info")

	d := doctor.Doctor{
		Name:            cleanText(e.ChildText("h1.entry-title")),
		BMDCRegNo:       extractBMDCRegNo(doctorInfo),
		Degrees:         cleanText(doctorInfo.Find(`li[title="Degree"]`).Text()),
		ExperienceYears: parseExperienceYears(doctorInfo.Find(`li[title="Experiences"]`).Text()),
		Specialties:     cleanText(doctorInfo.Find(`li[title="Specialty"]`).Text()),
		Designation:     cleanText(doctorInfo.Find(`li[title="Designation"]`).Text()),
		Workplace:       cleanText(doctorInfo.Find(`li[title="Workplace"]`).Text()),
		ImageURL:        e.ChildAttr(".photo img", "src"),
		ProfileURL:      e.Request.URL.String(),
	}

	e.DOM.Find(".entry-content h2").Each(func(_ int, h2 *goquery.Selection) {
		heading := cleanText(h2.Text())
		if !chamberHeadingRegex.MatchString(heading) {
			return // skip the Bangla SEO-content headings — not chamber data
		}

		p := h2.Next()
		if !p.Is("p") {
			return
		}

		chamber := parseChamberParagraph(p)
		if chamber.Name != "" || chamber.Address != "" {
			d.Chambers = append(d.Chambers, chamber)
		}
	})

	return d
}
