package aggregator

import "deepsearch/internal/model"

type Aggregator struct{}

func (a *Aggregator) Process(q model.Query, raw []model.Result) []model.Result {
	deduped := Dedup(raw)
	ranked := Rank(q, deduped)
	if q.MaxResults > 0 && len(ranked) > q.MaxResults {
		ranked = ranked[:q.MaxResults]
	}
	return ranked
}
