package web

import (
	"context"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	infrahttp "hybridsearch/internal/infra/http"
)

type RemoteResult struct {
	ID       string
	Title    string
	URL      string
	Snippet  string
	Score    float64
	Metadata map[string]string
}

type Client interface {
	Search(ctx context.Context, query string, limit int) ([]RemoteResult, error)
}

type DuckDuckGoLiteClient struct {
	httpClient *http.Client
	baseURL    string
}

func NewDuckDuckGoLiteClient(timeout time.Duration) *DuckDuckGoLiteClient {
	return &DuckDuckGoLiteClient{
		httpClient: infrahttp.NewClient(timeout),
		baseURL:    "https://lite.duckduckgo.com/lite/",
	}
}

func (c *DuckDuckGoLiteClient) Search(ctx context.Context, query string, limit int) ([]RemoteResult, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return []RemoteResult{}, nil
	}
	if limit <= 0 {
		limit = 10
	}

	u, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, fmt.Errorf("parse duckduckgo lite url: %w", err)
	}

	q := u.Query()
	q.Set("q", query)
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("create duckduckgo lite request: %w", err)
	}

	req.Header.Set("User-Agent", "hybridsearch/0.1 (+terminal deep search)")
	req.Header.Set("Accept", "text/html,application/xhtml+xml")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("duckduckgo lite request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("duckduckgo lite returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read duckduckgo lite response: %w", err)
	}

	return parseDuckDuckGoLiteHTML(string(body), limit), nil
}

var (
	anchorRegex = regexp.MustCompile(`(?is)<a[^>]+href="([^"]+)"[^>]*>(.*?)</a>`)
	tagRegex    = regexp.MustCompile(`(?is)<[^>]*>`)
)

func parseDuckDuckGoLiteHTML(raw string, limit int) []RemoteResult {
	matches := anchorRegex.FindAllStringSubmatch(raw, -1)
	if len(matches) == 0 {
		return []RemoteResult{}
	}

	results := make([]RemoteResult, 0, limit)
	seen := make(map[string]struct{})

	for _, m := range matches {
		if len(m) < 3 {
			continue
		}

		href := strings.TrimSpace(html.UnescapeString(m[1]))
		title := cleanHTMLText(m[2])

		target := extractRealTarget(href)
		if title == "" || target == "" {
			continue
		}

		if !strings.HasPrefix(target, "http://") && !strings.HasPrefix(target, "https://") {
			continue
		}

		if _, ok := seen[target]; ok {
			continue
		}
		seen[target] = struct{}{}

		results = append(results, RemoteResult{
			ID:      target,
			Title:   title,
			URL:     target,
			Snippet: "",
			Score:   1,
			Metadata: map[string]string{
				"provider": "duckduckgo-lite",
			},
		})

		if len(results) >= limit {
			break
		}
	}

	return results
}

func extractRealTarget(href string) string {
	if href == "" {
		return ""
	}

	parsed, err := url.Parse(href)
	if err == nil {
		if uddg := parsed.Query().Get("uddg"); uddg != "" {
			decoded, decodeErr := url.QueryUnescape(uddg)
			if decodeErr == nil {
				return decoded
			}
			return uddg
		}
	}

	return href
}

func cleanHTMLText(s string) string {
	s = tagRegex.ReplaceAllString(s, "")
	s = html.UnescapeString(s)
	s = strings.Join(strings.Fields(s), " ")
	return strings.TrimSpace(s)
}

// Optional helper if you want JSON-based fallback later.
func prettyJSON(v any) string {
	b, _ := json.Marshal(v)
	return string(b)
}
