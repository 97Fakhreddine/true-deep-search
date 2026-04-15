// Create new file: internal/provider/api/youtube_provider.go
package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	infrahttp "hybridsearch/internal/infra/http"
	hybridsearch "hybridsearch/pkg/hybridsearch"
)

type YouTubeProvider struct {
	client  *http.Client
	baseURL string
}

func NewYouTubeProvider(timeout time.Duration) *YouTubeProvider {
	return &YouTubeProvider{
		client:  infrahttp.NewClient(timeout),
		baseURL: "https://www.youtube.com/results",
	}
}

func (p *YouTubeProvider) Name() string {
	return "youtube"
}

func (p *YouTubeProvider) Kind() hybridsearch.ResultKind {
	return hybridsearch.ResultKindAPI
}

func (p *YouTubeProvider) Search(ctx context.Context, req hybridsearch.SearchRequest) ([]hybridsearch.SearchResult, error) {
	query := strings.TrimSpace(req.Query)
	if query == "" {
		return []hybridsearch.SearchResult{}, nil
	}

	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}
	if limit > 20 {
		limit = 20
	}

	u, err := url.Parse(p.baseURL)
	if err != nil {
		return nil, fmt.Errorf("parse youtube url: %w", err)
	}

	q := u.Query()
	q.Set("search_query", query)
	q.Set("pbj", "1")
	u.RawQuery = q.Encode()

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("create youtube request: %w", err)
	}

	httpReq.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0 Safari/537.36")
	httpReq.Header.Set("Accept", "application/json,text/plain,*/*")
	httpReq.Header.Set("X-YouTube-Client-Name", "1")
	httpReq.Header.Set("X-YouTube-Client-Version", "2.20241010.00.00")

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("youtube request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("youtube returned status %d", resp.StatusCode)
	}

	var payload any
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("decode youtube response: %w", err)
	}

	results := extractYouTubeResults(payload, limit)
	if len(results) == 0 {
		return []hybridsearch.SearchResult{}, nil
	}

	return results, nil
}

func extractYouTubeResults(payload any, limit int) []hybridsearch.SearchResult {
	var out []hybridsearch.SearchResult

	walkJSON(payload, func(node map[string]any) bool {
		videoRenderer, ok := node["videoRenderer"].(map[string]any)
		if !ok {
			return false
		}

		result, ok := mapYouTubeVideo(videoRenderer)
		if !ok {
			return false
		}

		out = append(out, result)
		return len(out) >= limit
	})

	return dedupeYouTubeResults(out)
}

func mapYouTubeVideo(video map[string]any) (hybridsearch.SearchResult, bool) {
	videoID := getString(video, "videoId")
	if strings.TrimSpace(videoID) == "" {
		return hybridsearch.SearchResult{}, false
	}

	title := extractRunsText(video["title"])
	if title == "" {
		title = extractSimpleText(video["title"])
	}
	title = strings.TrimSpace(title)
	if title == "" {
		return hybridsearch.SearchResult{}, false
	}

	target := "https://www.youtube.com/watch?v=" + videoID

	snippetParts := make([]string, 0, 4)

	owner := extractRunsText(video["ownerText"])
	if owner == "" {
		owner = extractRunsText(video["longBylineText"])
	}
	if owner != "" {
		snippetParts = append(snippetParts, owner)
	}

	length := extractSimpleText(video["lengthText"])
	if length != "" {
		snippetParts = append(snippetParts, "Duration: "+length)
	}

	views := extractSimpleText(video["viewCountText"])
	if views == "" {
		views = extractSimpleText(video["shortViewCountText"])
	}
	if views != "" {
		snippetParts = append(snippetParts, views)
	}

	published := extractSimpleText(video["publishedTimeText"])
	if published != "" {
		snippetParts = append(snippetParts, published)
	}

	description := extractRunsText(video["detailedMetadataSnippets"])
	if description == "" {
		description = extractRunsText(video["descriptionSnippet"])
	}
	description = strings.TrimSpace(description)

	snippet := strings.Join(snippetParts, " • ")
	if description != "" {
		if snippet != "" {
			snippet += " • "
		}
		snippet += description
	}

	score := 3.0
	if views != "" {
		score += 1
	}
	if published != "" {
		score += 1
	}
	if owner != "" {
		score += 1
	}

	return hybridsearch.SearchResult{
		ID:      videoID,
		Title:   title,
		Target:  target,
		Snippet: snippet,
		Source:  "youtube",
		Kind:    hybridsearch.ResultKindAPI,
		Score:   score,
		Metadata: map[string]string{
			"provider":  "youtube",
			"video_id":  videoID,
			"channel":   owner,
			"duration":  length,
			"views":     views,
			"published": published,
		},
	}, true
}

func dedupeYouTubeResults(results []hybridsearch.SearchResult) []hybridsearch.SearchResult {
	if len(results) == 0 {
		return []hybridsearch.SearchResult{}
	}

	seen := make(map[string]struct{}, len(results))
	out := make([]hybridsearch.SearchResult, 0, len(results))

	for _, r := range results {
		key := strings.TrimSpace(r.ID)
		if key == "" {
			key = strings.TrimSpace(r.Target)
		}
		if key == "" {
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

func walkJSON(v any, visit func(map[string]any) bool) bool {
	switch node := v.(type) {
	case map[string]any:
		if visit(node) {
			return true
		}
		for _, child := range node {
			if walkJSON(child, visit) {
				return true
			}
		}
	case []any:
		for _, child := range node {
			if walkJSON(child, visit) {
				return true
			}
		}
	}
	return false
}

func extractRunsText(v any) string {
	switch node := v.(type) {
	case map[string]any:
		if runs, ok := node["runs"].([]any); ok {
			parts := make([]string, 0, len(runs))
			for _, item := range runs {
				if runMap, ok := item.(map[string]any); ok {
					text := getString(runMap, "text")
					if text != "" {
						parts = append(parts, text)
					}
				}
			}
			return strings.Join(parts, "")
		}
		if simple := getString(node, "simpleText"); simple != "" {
			return simple
		}
	case []any:
		parts := make([]string, 0, len(node))
		for _, item := range node {
			text := extractRunsText(item)
			if text != "" {
				parts = append(parts, text)
			}
		}
		return strings.Join(parts, " ")
	}
	return ""
}

func extractSimpleText(v any) string {
	if m, ok := v.(map[string]any); ok {
		if s := getString(m, "simpleText"); s != "" {
			return s
		}
	}
	return ""
}

func getString(m map[string]any, key string) string {
	raw, ok := m[key]
	if !ok || raw == nil {
		return ""
	}

	switch v := raw.(type) {
	case string:
		return v
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case json.Number:
		return v.String()
	default:
		return ""
	}
}
