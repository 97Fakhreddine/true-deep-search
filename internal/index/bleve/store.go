package bleve

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/blevesearch/bleve/v2/search"

	hybridsearch "hybridsearch/pkg/hybridsearch"
)

type storeDocument struct {
	ID      string `json:"id"`
	Path    string `json:"path"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func mapSearchResults(hits search.DocumentMatchCollection) []hybridsearch.SearchResult {
	if len(hits) == 0 {
		return []hybridsearch.SearchResult{}
	}

	out := make([]hybridsearch.SearchResult, 0, len(hits))

	for _, hit := range hits {
		title := fieldAsString(hit.Fields["title"])
		target := fieldAsString(hit.Fields["path"])
		content := fieldAsString(hit.Fields["content"])

		if strings.TrimSpace(title) == "" {
			title = fallbackTitle(target, hit.ID)
		}

		out = append(out, hybridsearch.SearchResult{
			ID:      hit.ID,
			Title:   title,
			Target:  target,
			Snippet: makeSnippet(content, 180),
			Source:  "local",
			Kind:    hybridsearch.ResultKindLocal,
			Score:   hit.Score,
			Metadata: map[string]string{
				"engine": "bleve",
			},
		})
	}

	return out
}

func fieldAsString(v any) string {
	if v == nil {
		return ""
	}

	switch t := v.(type) {
	case string:
		return t
	case []byte:
		return string(t)
	default:
		return fmt.Sprint(v)
	}
}

func fallbackTitle(path string, fallback string) string {
	path = strings.TrimSpace(path)
	if path == "" {
		return fallback
	}

	return filepath.Base(path)
}

func makeSnippet(content string, max int) string {
	content = strings.TrimSpace(content)
	if content == "" {
		return ""
	}

	content = strings.ReplaceAll(content, "\n", " ")
	content = strings.ReplaceAll(content, "\r", " ")
	content = strings.Join(strings.Fields(content), " ")

	if max <= 0 || len(content) <= max {
		return content
	}

	if max <= 3 {
		return content[:max]
	}

	return content[:max-3] + "..."
}
