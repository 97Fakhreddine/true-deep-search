package aggregate

import hybridsearch "hybridsearch/pkg/hybridsearch"

type Merger struct{}

func NewMerger() *Merger {
	return &Merger{}
}

func (m *Merger) Merge(inputs map[string][]hybridsearch.SearchResult) []hybridsearch.SearchResult {
	if len(inputs) == 0 {
		return []hybridsearch.SearchResult{}
	}

	total := 0
	for _, results := range inputs {
		total += len(results)
	}

	merged := make([]hybridsearch.SearchResult, 0, total)

	for providerName, results := range inputs {
		for _, r := range results {
			if r.Source == "" {
				r.Source = providerName
			}
			merged = append(merged, r)
		}
	}

	return merged
}
