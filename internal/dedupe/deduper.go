package dedupe

import hybridsearch "hybridsearch/pkg/hybridsearch"

type Deduper struct{}

func New() *Deduper {
	return &Deduper{}
}

func (d *Deduper) Dedupe(results []hybridsearch.SearchResult) []hybridsearch.SearchResult {
	if len(results) == 0 {
		return []hybridsearch.SearchResult{}
	}

	seen := make(map[string]struct{}, len(results))
	out := make([]hybridsearch.SearchResult, 0, len(results))

	for _, r := range results {
		key := fingerprint(r)

		if key == "" {
			out = append(out, r)
			continue
		}

		if _, ok := seen[key]; ok {
			continue
		}

		seen[key] = struct{}{}
		out = append(out, r)
	}

	return out
}
