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

type StackExchangeProvider struct {
	client  *http.Client
	baseURL string
	site    string
}

func NewStackExchangeProvider(timeout time.Duration) *StackExchangeProvider {
	return &StackExchangeProvider{
		client:  infrahttp.NewClient(timeout),
		baseURL: "https://api.stackexchange.com/2.3/search/advanced",
		site:    "stackoverflow",
	}
}

func (p *StackExchangeProvider) Name() string {
	return "stackexchange"
}

func (p *StackExchangeProvider) Kind() hybridsearch.ResultKind {
	return hybridsearch.ResultKindAPI
}

func (p *StackExchangeProvider) Search(ctx context.Context, req hybridsearch.SearchRequest) ([]hybridsearch.SearchResult, error) {
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
		return nil, fmt.Errorf("parse stackexchange url: %w", err)
	}

	q := u.Query()
	q.Set("order", "desc")
	q.Set("sort", "relevance")
	q.Set("q", query)
	q.Set("site", p.site)
	q.Set("pagesize", strconv.Itoa(limit))
	u.RawQuery = q.Encode()

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("create stackexchange request: %w", err)
	}

	httpReq.Header.Set("User-Agent", "hybridsearch/0.1 (+terminal deep search)")

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("stackexchange request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("stackexchange returned status %d", resp.StatusCode)
	}

	var payload stackExchangeResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("decode stackexchange response: %w", err)
	}

	return mapStackExchangeResults(payload), nil
}

type stackExchangeResponse struct {
	Items []stackExchangeItem `json:"items"`
}

type stackExchangeItem struct {
	QuestionID int64    `json:"question_id"`
	Title      string   `json:"title"`
	Link       string   `json:"link"`
	IsAnswered bool     `json:"is_answered"`
	Score      int      `json:"score"`
	Tags       []string `json:"tags"`
}

func mapStackExchangeResults(payload stackExchangeResponse) []hybridsearch.SearchResult {
	if len(payload.Items) == 0 {
		return []hybridsearch.SearchResult{}
	}

	out := make([]hybridsearch.SearchResult, 0, len(payload.Items))

	for _, item := range payload.Items {
		title := strings.TrimSpace(item.Title)
		target := strings.TrimSpace(item.Link)

		if title == "" || target == "" {
			continue
		}

		snippet := buildStackExchangeSnippet(item)

		out = append(out, hybridsearch.SearchResult{
			ID:      strconv.FormatInt(item.QuestionID, 10),
			Title:   htmlUnescape(title),
			Target:  target,
			Snippet: snippet,
			Source:  "stackexchange",
			Kind:    hybridsearch.ResultKindAPI,
			Score:   float64(item.Score) + boolToScore(item.IsAnswered),
			Metadata: map[string]string{
				"provider": "stackexchange",
				"answered": strconv.FormatBool(item.IsAnswered),
			},
		})
	}

	return out
}

func buildStackExchangeSnippet(item stackExchangeItem) string {
	parts := []string{}

	if len(item.Tags) > 0 {
		parts = append(parts, "Tags: "+strings.Join(item.Tags[:min(3, len(item.Tags))], ", "))
	}

	if item.IsAnswered {
		parts = append(parts, "Answered")
	} else {
		parts = append(parts, "Unanswered")
	}

	parts = append(parts, fmt.Sprintf("Score: %d", item.Score))

	return strings.Join(parts, " • ")
}

func boolToScore(v bool) float64 {
	if v {
		return 5
	}
	return 0
}

func htmlUnescape(s string) string {
	replacer := strings.NewReplacer(
		"&amp;", "&",
		"&lt;", "<",
		"&gt;", ">",
		"&quot;", "\"",
		"&#39;", "'",
	)
	return replacer.Replace(s)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
