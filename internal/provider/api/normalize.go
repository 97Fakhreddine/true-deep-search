package api

import (
	"strings"

	hybridsearch "hybridsearch/pkg/hybridsearch"
)

type APIResult struct {
	ID       string
	Title    string
	Target   string
	Snippet  string
	Score    float64
	Metadata map[string]string
}

func normalizeAPIResults(results []APIResult, source string) []hybridsearch.SearchResult {
	if len(results) == 0 {
		return []hybridsearch.SearchResult{}
	}

	out := make([]hybridsearch.SearchResult, 0, len(results))

	for _, r := range results {
		meta := r.Metadata
		if meta == nil {
			meta = map[string]string{}
		}

		out = append(out, hybridsearch.SearchResult{
			ID:       strings.TrimSpace(r.ID),
			Title:    strings.TrimSpace(r.Title),
			Target:   strings.TrimSpace(r.Target),
			Snippet:  strings.TrimSpace(r.Snippet),
			Source:   source,
			Kind:     hybridsearch.ResultKindAPI,
			Score:    r.Score,
			Metadata: meta,
		})
	}

	return out
}
