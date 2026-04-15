package search

type Aggregator interface {
	Merge(inputs map[string][]SearchResult) []SearchResult
}

type Deduper interface {
	Dedupe(results []SearchResult) []SearchResult
}

type Ranker interface {
	Rank(query string, results []SearchResult) []SearchResult
}
