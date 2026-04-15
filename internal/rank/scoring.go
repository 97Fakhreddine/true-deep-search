// Replace file: internal/rank/scoring.go
package rank

import (
	"strings"

	hybridsearch "hybridsearch/pkg/hybridsearch"
)

func computeScore(query string, r hybridsearch.SearchResult) float64 {
	score := r.Score

	q := normalize(query)
	if q == "" {
		return score
	}

	title := normalize(r.Title)
	snippet := normalize(r.Snippet)
	target := normalize(r.Target)
	source := normalize(r.Source)

	terms := splitTerms(q)

	// Strong exact title boost
	if title == q {
		score += 120
	}

	// Partial title match
	if strings.Contains(title, q) {
		score += 60
	}

	// Snippet match
	if strings.Contains(snippet, q) {
		score += 25
	}

	// Target/path/url match
	if strings.Contains(target, q) {
		score += 15
	}

	// Term coverage
	score += termCoverageBoost(terms, title, 12)
	score += termCoverageBoost(terms, snippet, 6)
	score += termCoverageBoost(terms, target, 4)

	// Base source weighting
	score += sourceBaseBoost(source)

	// Intent-aware source boosts
	score += intentBoost(q, source)

	// Local slight boost
	if r.Kind == hybridsearch.ResultKindLocal {
		score += 6
	}

	// Metadata-aware boosts
	score += metadataBoost(r)

	return score
}

func normalize(s string) string {
	return strings.TrimSpace(strings.ToLower(s))
}

func splitTerms(q string) []string {
	raw := strings.Fields(q)
	out := make([]string, 0, len(raw))

	seen := make(map[string]struct{}, len(raw))
	for _, term := range raw {
		term = normalize(term)
		if len(term) < 2 {
			continue
		}
		if _, ok := seen[term]; ok {
			continue
		}
		seen[term] = struct{}{}
		out = append(out, term)
	}

	return out
}

func termCoverageBoost(terms []string, text string, perTerm float64) float64 {
	if len(terms) == 0 || text == "" {
		return 0
	}

	var score float64
	for _, term := range terms {
		if strings.Contains(text, term) {
			score += perTerm
		}
	}
	return score
}

func sourceBaseBoost(source string) float64 {
	switch source {
	case "wikipedia":
		return 18
	case "github":
		return 20
	case "stackexchange":
		return 22
	case "reddit":
		return 12
	case "youtube":
		return 14
	case "web":
		return 10
	case "local":
		return 8
	default:
		return 0
	}
}

func intentBoost(query string, source string) float64 {
	var boost float64

	if isTechnicalQuery(query) {
		switch source {
		case "github":
			boost += 30
		case "stackexchange":
			boost += 35
		case "reddit":
			boost += 8
		case "youtube":
			boost += 10
		case "wikipedia":
			boost += 4
		}
	}

	if isTutorialQuery(query) {
		switch source {
		case "youtube":
			boost += 35
		case "github":
			boost += 15
		case "stackexchange":
			boost += 10
		case "reddit":
			boost += 6
		}
	}

	if isKnowledgeQuery(query) {
		switch source {
		case "wikipedia":
			boost += 35
		case "web":
			boost += 12
		case "youtube":
			boost += 5
		}
	}

	if isCommunityQuery(query) {
		switch source {
		case "reddit":
			boost += 30
		case "stackexchange":
			boost += 15
		case "youtube":
			boost += 8
		}
	}

	if isVideoQuery(query) {
		switch source {
		case "youtube":
			boost += 40
		case "web":
			boost += 5
		}
	}

	return boost
}

func isTechnicalQuery(q string) bool {
	keywords := []string{
		"golang", "go", "javascript", "typescript", "react", "angular", "vue",
		"node", "nestjs", "api", "sql", "postgres", "mongodb", "docker",
		"kubernetes", "bug", "error", "exception", "stack trace", "library",
		"framework", "code", "coding", "programming", "dev", "developer",
		"algorithm", "binary", "compile", "build", "test", "ci", "cd",
	}
	return containsAny(q, keywords)
}

func isTutorialQuery(q string) bool {
	keywords := []string{
		"tutorial", "course", "learn", "how to", "guide", "walkthrough",
		"step by step", "lesson", "crash course", "explained", "beginner",
	}
	return containsAny(q, keywords)
}

func isKnowledgeQuery(q string) bool {
	keywords := []string{
		"what is", "who is", "history", "meaning", "definition", "overview",
		"explain", "encyclopedia", "concept", "theory",
	}
	return containsAny(q, keywords)
}

func isCommunityQuery(q string) bool {
	keywords := []string{
		"opinion", "review", "experience", "discussion", "community",
		"reddit", "best practice", "thoughts", "compare", "vs",
	}
	return containsAny(q, keywords)
}

func isVideoQuery(q string) bool {
	keywords := []string{
		"video", "youtube", "watch", "stream", "recording", "talk", "conference",
	}
	return containsAny(q, keywords)
}

func containsAny(q string, keywords []string) bool {
	for _, k := range keywords {
		if strings.Contains(q, k) {
			return true
		}
	}
	return false
}

func metadataBoost(r hybridsearch.SearchResult) float64 {
	if len(r.Metadata) == 0 {
		return 0
	}

	var score float64

	switch normalize(r.Source) {
	case "github":
		if lang := normalize(r.Metadata["language"]); lang != "" {
			score += 3
		}

	case "stackexchange":
		if normalize(r.Metadata["answered"]) == "true" {
			score += 10
		}

	case "reddit":
		if comments := normalize(r.Metadata["comments"]); comments != "" {
			score += 2
		}

	case "youtube":
		if channel := normalize(r.Metadata["channel"]); channel != "" {
			score += 3
		}
		if duration := normalize(r.Metadata["duration"]); duration != "" {
			score += 2
		}
	}

	return score
}
