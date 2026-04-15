package local

import (
	"strings"

	hybridsearch "hybridsearch/pkg/hybridsearch"
)

func normalizeResults(results []hybridsearch.SearchResult) []hybridsearch.SearchResult {
	if len(results) == 0 {
		return []hybridsearch.SearchResult{}
	}

	out := make([]hybridsearch.SearchResult, 0, len(results))

	for _, r := range results {
		r.Source = "local"
		r.Kind = hybridsearch.ResultKindLocal
		r.Title = strings.TrimSpace(r.Title)
		r.Target = strings.TrimSpace(r.Target)
		r.Snippet = strings.TrimSpace(r.Snippet)

		if r.Metadata == nil {
			r.Metadata = map[string]string{}
		}

		out = append(out, r)
	}

	return out
}
