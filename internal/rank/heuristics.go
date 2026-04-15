package rank

import (
	"strings"

	hybridsearch "hybridsearch/pkg/hybridsearch"
)

func applyHeuristics(results []hybridsearch.SearchResult) []hybridsearch.SearchResult {
	for i := range results {
		r := &results[i]

		title := strings.TrimSpace(r.Title)
		target := strings.TrimSpace(r.Target)
		snippet := strings.TrimSpace(r.Snippet)
		source := normalize(r.Source)

		// Penalize weak/empty fields
		if title == "" {
			r.Score -= 15
		}
		if target == "" {
			r.Score -= 15
		}
		if snippet == "" {
			r.Score -= 4
		}

		// Reward richer snippets a bit
		if len(snippet) >= 40 {
			r.Score += 3
		}
		if len(snippet) >= 120 {
			r.Score += 2
		}

		// Small quality boost for https targets
		if strings.HasPrefix(strings.ToLower(target), "https://") {
			r.Score += 2
		}

		// Slight penalty for obvious tracking/search redirect style links
		if looksLikeTrackingOrRedirect(target) {
			r.Score -= 6
		}

		// Prefer canonical-looking result pages
		if looksLikeSearchPage(target) {
			r.Score -= 10
		}

		// Source-specific heuristics
		switch source {
		case "github":
			if strings.Contains(strings.ToLower(target), "/issues") {
				r.Score -= 2
			}
			if strings.Contains(strings.ToLower(target), "/pull/") {
				r.Score -= 2
			}

		case "reddit":
			if normalize(r.Metadata["nsfw"]) == "true" {
				r.Score -= 20
			}

		case "youtube":
			if strings.Contains(strings.ToLower(title), "shorts") {
				r.Score -= 2
			}

		case "wikipedia":
			// Wikipedia pages with decent snippet are usually useful
			if len(snippet) > 0 {
				r.Score += 4
			}

		case "stackexchange":
			if normalize(r.Metadata["answered"]) == "true" {
				r.Score += 6
			}
		}

		// Never allow negative scores
		if r.Score < 0 {
			r.Score = 0
		}
	}

	return results
}

func looksLikeTrackingOrRedirect(target string) bool {
	t := strings.ToLower(strings.TrimSpace(target))
	if t == "" {
		return false
	}

	badParts := []string{
		"utm_",
		"fbclid=",
		"gclid=",
		"redirect=",
		"redirect_uri=",
		"ref=",
		"ref_src=",
	}

	for _, part := range badParts {
		if strings.Contains(t, part) {
			return true
		}
	}

	return false
}

func looksLikeSearchPage(target string) bool {
	t := strings.ToLower(strings.TrimSpace(target))
	if t == "" {
		return false
	}

	searchPatterns := []string{
		"/search?",
		"/search/",
		"search_query=",
		"?q=",
		"&q=",
		"?query=",
		"&query=",
		"?search=",
		"&search=",
	}

	for _, p := range searchPatterns {
		if strings.Contains(t, p) {
			return true
		}
	}

	return false
}
