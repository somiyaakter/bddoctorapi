package taxonomy

import (
	"regexp"
	"strings"
)

var locationAliases = map[string][]string{
	"chittagong": {"chittagong", "chattogram", "ctg"},
	"cumilla":    {"comilla", "cumilla"},
	"bogura":     {"bogra", "bogura"},
	"barisal":    {"barisal", "barishal"},
}

func slugFromLocationURL(locationURL string) string {
	trimmed := strings.TrimSuffix(locationURL, "/")
	idx := strings.LastIndex(trimmed, "doctors-")
	if idx == -1 {
		return ""
	}
	return strings.ToLower(trimmed[idx+len("doctors-"):])
}

// MatchLocations returns every Location whose name (or known alias)
// appears in any of the given chamber address strings.
func MatchLocations(locations []Location, addresses []string) []Location {
	haystack := strings.ToLower(strings.Join(addresses, " | "))

	var matched []Location
	for _, loc := range locations {
		names := []string{strings.ToLower(loc.Name)}
		if aliases, ok := locationAliases[slugFromLocationURL(loc.URL)]; ok {
			names = aliases
		}
		for _, name := range names {
			if name != "" && strings.Contains(haystack, name) {
				matched = append(matched, loc)
				break
			}
		}
	}
	return matched
}

var nonLetterRegex = regexp.MustCompile(`[^a-z]+`)

func specialtyStem(word string) string {
	word = nonLetterRegex.ReplaceAllString(strings.ToLower(word), "")
	const stemLen = 6
	if len(word) <= stemLen {
		return word
	}
	return word[:stemLen]
}

func specialtySlugWords(specialtyURL, locationSlug string) []string {
	trimmed := strings.TrimSuffix(specialtyURL, "/")
	idx := strings.LastIndex(trimmed, "/")
	if idx == -1 {
		return nil
	}
	slug := strings.TrimSuffix(trimmed[idx+1:], "-"+locationSlug)
	if slug == "" {
		return nil
	}
	return strings.Split(slug, "-")
}

// MatchSpecialties returns every Specialty whose slug-derived word
// stems ALL appear in specialtyText (the doctor's own scraped
// "Specialty" field). locationURLByID maps each specialty's
// LocationID to that location's URL, needed to strip the
// "-{locationSlug}" suffix off the specialty's own URL slug.
func MatchSpecialties(specialties []Specialty, locationURLByID map[int64]string, specialtyText string) []Specialty {
	haystack := strings.ToLower(specialtyText)

	var matched []Specialty
	for _, sp := range specialties {
		locSlug := slugFromLocationURL(locationURLByID[sp.LocationID])
		words := specialtySlugWords(sp.URL, locSlug)
		if len(words) == 0 {
			continue
		}

		allPresent := true
		for _, w := range words {
			stem := specialtyStem(w)
			if stem == "" || !strings.Contains(haystack, stem) {
				allPresent = false
				break
			}
		}
		if allPresent {
			matched = append(matched, sp)
		}
	}
	return matched
}