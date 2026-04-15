package rank

import (
	"sort"

	hybridsearch "hybridsearch/pkg/hybridsearch"
)

type Ranker struct{}

func New() *Ranker {
	return &Ranker{}
}

func (r *Ranker) Rank(query string, results []hybridsearch.SearchResult) []hybridsearch.SearchResult {
	if len(results) == 0 {
		return []hybridsearch.SearchResult{}
	}

	for i := range results {
		results[i].Score = computeScore(query, results[i])
	}

	results = applyHeuristics(results)

	sort.SliceStable(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	return results
}
