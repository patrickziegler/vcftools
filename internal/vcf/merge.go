package vcf

import (
	"sort"
	"strings"
)

func telKeys(c *Card) []string {
	var out []string
	for _, p := range c.All("TEL") {
		if v := strings.TrimSpace(p.Value); v != "" {
			out = append(out, v)
		}
	}
	return out
}

func nameKey(c *Card) string {
	full, _, _ := c.Names()
	return strings.ToLower(strings.Join(strings.Fields(full), " "))
}

// Merge removes duplicate cards without rewriting any property. Cards are
// duplicates when they have an identical TEL value or the same
// case-insensitive name (raw values are compared, so run cleanvcf first to
// get reliable matching). TEL lines of a duplicate that the kept card lacks
// are appended unchanged to the kept card so no phone number is lost.
func Merge(cards []*Card) []*Card {
	var out []*Card
	byPhone := map[string]*Card{}
	byName := map[string]*Card{}
	for _, c := range cards {
		key := nameKey(c)
		var match *Card
		for _, p := range telKeys(c) {
			if m, ok := byPhone[p]; ok {
				match = m
				break
			}
		}
		if match == nil && key != "" {
			match = byName[key]
		}
		if match == nil {
			match = c
			out = append(out, c)
		} else {
			have := map[string]bool{}
			for _, p := range telKeys(match) {
				have[p] = true
			}
			for _, p := range c.All("TEL") {
				v := strings.TrimSpace(p.Value)
				if v != "" && !have[v] {
					have[v] = true
					match.Props = append(match.Props, p)
				}
			}
		}
		for _, p := range telKeys(match) {
			byPhone[p] = match
		}
		if key != "" {
			byName[key] = match
		}
	}
	return out
}

// SortByLastName orders cards by family name, then given name, then full name.
func SortByLastName(cards []*Card) {
	type k struct{ full, family, given string }
	keys := make(map[*Card]k, len(cards))
	for _, c := range cards {
		f, fa, g := c.Names()
		keys[c] = k{strings.ToLower(f), strings.ToLower(fa), strings.ToLower(g)}
	}
	sort.SliceStable(cards, func(i, j int) bool {
		a, b := keys[cards[i]], keys[cards[j]]
		if a.family != b.family {
			return a.family < b.family
		}
		if a.given != b.given {
			return a.given < b.given
		}
		return a.full < b.full
	})
}
