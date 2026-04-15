package aggregator

import (
	"deepsearch/internal/model"
	"sort"
	"strings"
)

func Rank(q model.Query, results []model.Result) []model.Result {
	terms := strings.Fields(strings.ToLower(q.Text))
	for i := range results {
		results[i].Score = score(results[i], terms)
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})
	return results
}

func score(r model.Result, terms []string) float64 {
	base := r.Score
	text := strings.ToLower(r.Title + " " + r.Snippet)
	for _, t := range terms {
		if strings.Contains(r.Title, t) {
			base += 2.0
		} else if strings.Contains(text, t) {
			base += 0.5
		}
	}
	return base
}
