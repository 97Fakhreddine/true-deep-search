package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	infrahttp "hybridsearch/internal/infra/http"
	hybridsearch "hybridsearch/pkg/hybridsearch"
)

type WikipediaProvider struct {
	client  *http.Client
	baseURL string
}

func NewWikipediaProvider(timeout time.Duration) *WikipediaProvider {
	return &WikipediaProvider{
		client:  infrahttp.NewClient(timeout),
		baseURL: "https://en.wikipedia.org/w/api.php",
	}
}

func (p *WikipediaProvider) Name() string {
	return "wikipedia"
}

func (p *WikipediaProvider) Kind() hybridsearch.ResultKind {
	return hybridsearch.ResultKindAPI
}

func (p *WikipediaProvider) Search(ctx context.Context, req hybridsearch.SearchRequest) ([]hybridsearch.SearchResult, error) {
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
		return nil, fmt.Errorf("parse wikipedia url: %w", err)
	}

	q := u.Query()
	q.Set("action", "opensearch")
	q.Set("search", query)
	q.Set("limit", fmt.Sprintf("%d", limit))
	q.Set("namespace", "0")
	q.Set("format", "json")
	q.Set("origin", "*")
	u.RawQuery = q.Encode()

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("create wikipedia request: %w", err)
	}

	httpReq.Header.Set("User-Agent", "hybridsearch/0.1 (+terminal deep search)")

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("wikipedia request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("wikipedia returned status %d", resp.StatusCode)
	}

	var payload []any
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("decode wikipedia response: %w", err)
	}

	return mapWikipediaResponse(payload), nil
}

func mapWikipediaResponse(payload []any) []hybridsearch.SearchResult {
	if len(payload) < 4 {
		return []hybridsearch.SearchResult{}
	}

	titles := toStringSlice(payload[1])
	descriptions := toStringSlice(payload[2])
	links := toStringSlice(payload[3])

	size := minInt(len(titles), len(descriptions), len(links))
	if size == 0 {
		return []hybridsearch.SearchResult{}
	}

	out := make([]hybridsearch.SearchResult, 0, size)

	for i := 0; i < size; i++ {
		title := strings.TrimSpace(titles[i])
		target := strings.TrimSpace(links[i])
		snippet := strings.TrimSpace(descriptions[i])

		if title == "" || target == "" {
			continue
		}

		out = append(out, hybridsearch.SearchResult{
			ID:      target,
			Title:   title,
			Target:  target,
			Snippet: snippet,
			Source:  "wikipedia",
			Kind:    hybridsearch.ResultKindAPI,
			Score:   2,
			Metadata: map[string]string{
				"provider": "wikipedia",
			},
		})
	}

	return out
}

func toStringSlice(v any) []string {
	arr, ok := v.([]any)
	if !ok {
		return []string{}
	}

	out := make([]string, 0, len(arr))
	for _, item := range arr {
		s, ok := item.(string)
		if !ok {
			out = append(out, "")
			continue
		}
		out = append(out, s)
	}
	return out
}

func minInt(vals ...int) int {
	if len(vals) == 0 {
		return 0
	}
	min := vals[0]
	for _, v := range vals[1:] {
		if v < min {
			min = v
		}
	}
	return min
}
