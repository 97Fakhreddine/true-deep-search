package aggregator

import (
	"crypto/sha256"
	"deepsearch/internal/model"
	"fmt"
)

func Dedup(results []model.Result) []model.Result {
	seen := make(map[string]bool)
	out := results[:0]
	for _, r := range results {
		key := fingerprint(r)
		if !seen[key] {
			seen[key] = true
			out = append(out, r)
		}
	}
	return out
}

func fingerprint(r model.Result) string {
	if r.URL != "" {
		return r.URL
	}
	h := sha256.Sum256([]byte(r.Title + r.Snippet))
	return fmt.Sprintf("%x", h[:8])
}
