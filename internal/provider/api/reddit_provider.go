// Create new file: internal/provider/api/reddit_provider.go
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

type RedditProvider struct {
	client  *http.Client
	baseURL string
}

func NewRedditProvider(timeout time.Duration) *RedditProvider {
	return &RedditProvider{
		client:  infrahttp.NewClient(timeout),
		baseURL: "https://www.reddit.com/search.json",
	}
}

func (p *RedditProvider) Name() string {
	return "reddit"
}

func (p *RedditProvider) Kind() hybridsearch.ResultKind {
	return hybridsearch.ResultKindAPI
}

func (p *RedditProvider) Search(ctx context.Context, req hybridsearch.SearchRequest) ([]hybridsearch.SearchResult, error) {
	query := strings.TrimSpace(req.Query)
	if query == "" {
		return []hybridsearch.SearchResult{}, nil
	}

	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}
	if limit > 25 {
		limit = 25
	}

	u, err := url.Parse(p.baseURL)
	if err != nil {
		return nil, fmt.Errorf("parse reddit url: %w", err)
	}

	q := u.Query()
	q.Set("q", query)
	q.Set("limit", strconv.Itoa(limit))
	q.Set("sort", "relevance")
	q.Set("type", "link")
	q.Set("raw_json", "1")
	u.RawQuery = q.Encode()

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("create reddit request: %w", err)
	}

	httpReq.Header.Set("User-Agent", "hybridsearch/0.1 (+terminal deep search)")
	httpReq.Header.Set("Accept", "application/json")

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("reddit request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("reddit returned status %d", resp.StatusCode)
	}

	var payload redditSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("decode reddit response: %w", err)
	}

	return mapRedditResults(payload), nil
}

type redditSearchResponse struct {
	Data redditListingData `json:"data"`
}

type redditListingData struct {
	Children []redditChild `json:"children"`
}

type redditChild struct {
	Data redditPost `json:"data"`
}

type redditPost struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	SelfText    string  `json:"selftext"`
	Permalink   string  `json:"permalink"`
	Subreddit   string  `json:"subreddit"`
	Author      string  `json:"author"`
	Score       int     `json:"score"`
	NumComments int     `json:"num_comments"`
	CreatedUTC  float64 `json:"created_utc"`
	IsSelf      bool    `json:"is_self"`
	URL         string  `json:"url"`
	Over18      bool    `json:"over_18"`
}

func mapRedditResults(payload redditSearchResponse) []hybridsearch.SearchResult {
	if len(payload.Data.Children) == 0 {
		return []hybridsearch.SearchResult{}
	}

	out := make([]hybridsearch.SearchResult, 0, len(payload.Data.Children))

	for _, child := range payload.Data.Children {
		post := child.Data

		title := strings.TrimSpace(post.Title)
		if title == "" {
			continue
		}

		target := buildRedditTarget(post)
		if target == "" {
			continue
		}

		snippet := buildRedditSnippet(post)

		meta := map[string]string{
			"provider":     "reddit",
			"subreddit":    post.Subreddit,
			"author":       post.Author,
			"comments":     strconv.Itoa(post.NumComments),
			"nsfw":         strconv.FormatBool(post.Over18),
			"is_self_post": strconv.FormatBool(post.IsSelf),
		}

		if post.CreatedUTC > 0 {
			meta["created_utc"] = fmt.Sprintf("%.0f", post.CreatedUTC)
		}

		out = append(out, hybridsearch.SearchResult{
			ID:       strings.TrimSpace(post.ID),
			Title:    title,
			Target:   target,
			Snippet:  snippet,
			Source:   "reddit",
			Kind:     hybridsearch.ResultKindAPI,
			Score:    float64(post.Score) + commentBoost(post.NumComments),
			Metadata: meta,
		})
	}

	return out
}

func buildRedditTarget(post redditPost) string {
	permalink := strings.TrimSpace(post.Permalink)
	if permalink != "" {
		if strings.HasPrefix(permalink, "http://") || strings.HasPrefix(permalink, "https://") {
			return permalink
		}
		return "https://www.reddit.com" + permalink
	}

	if strings.TrimSpace(post.URL) != "" {
		return strings.TrimSpace(post.URL)
	}

	return ""
}

func buildRedditSnippet(post redditPost) string {
	parts := make([]string, 0, 4)

	if post.Subreddit != "" {
		parts = append(parts, "r/"+post.Subreddit)
	}

	if post.Author != "" {
		parts = append(parts, "u/"+post.Author)
	}

	parts = append(parts, fmt.Sprintf("Score: %d", post.Score))
	parts = append(parts, fmt.Sprintf("%d comments", post.NumComments))

	base := strings.Join(parts, " • ")

	body := makeCompactSnippet(post.SelfText, 160)
	if body == "" {
		return base
	}

	return base + " • " + body
}

func commentBoost(comments int) float64 {
	switch {
	case comments >= 500:
		return 8
	case comments >= 100:
		return 5
	case comments >= 20:
		return 3
	case comments > 0:
		return 1
	default:
		return 0
	}
}

func makeCompactSnippet(s string, max int) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}

	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\r", " ")
	s = strings.Join(strings.Fields(s), " ")

	if max <= 0 || len(s) <= max {
		return s
	}

	if max <= 3 {
		return s[:max]
	}

	return s[:max-3] + "..."
}
