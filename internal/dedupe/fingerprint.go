package dedupe

import (
	"strings"

	hybridsearch "hybridsearch/pkg/hybridsearch"
)

func fingerprint(r hybridsearch.SearchResult) string {
	target := normalize(r.Target)
	if target != "" {
		return target
	}

	title := normalize(r.Title)
	return title
}

func normalize(s string) string {
	return strings.TrimSpace(strings.ToLower(s))
}
