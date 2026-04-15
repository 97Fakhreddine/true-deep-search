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

type GitHubProvider struct {
	client  *http.Client
	baseURL string
	token   string
}

func NewGitHubProvider(timeout time.Duration, token string) *GitHubProvider {
	return &GitHubProvider{
		client:  infrahttp.NewClient(timeout),
		baseURL: "https://api.github.com/search/repositories",
		token:   strings.TrimSpace(token),
	}
}

func (p *GitHubProvider) Name() string {
	return "github"
}

func (p *GitHubProvider) Kind() hybridsearch.ResultKind {
	return hybridsearch.ResultKindAPI
}

func (p *GitHubProvider) Search(ctx context.Context, req hybridsearch.SearchRequest) ([]hybridsearch.SearchResult, error) {
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
		return nil, fmt.Errorf("parse github url: %w", err)
	}

	q := u.Query()
	q.Set("q", query)
	q.Set("sort", "best-match")
	q.Set("order", "desc")
	q.Set("per_page", strconv.Itoa(limit))
	u.RawQuery = q.Encode()

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("create github request: %w", err)
	}

	httpReq.Header.Set("User-Agent", "hybridsearch/0.1 (+terminal deep search)")
	httpReq.Header.Set("Accept", "application/vnd.github+json")
	if p.token != "" {
		httpReq.Header.Set("Authorization", "Bearer "+p.token)
	}

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("github request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("github returned status %d", resp.StatusCode)
	}

	var payload githubSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("decode github response: %w", err)
	}

	return mapGitHubResponse(payload), nil
}

type githubSearchResponse struct {
	Items []githubRepo `json:"items"`
}

type githubRepo struct {
	ID              int64  `json:"id"`
	FullName        string `json:"full_name"`
	HTMLURL         string `json:"html_url"`
	Description     string `json:"description"`
	Language        string `json:"language"`
	StargazersCount int64  `json:"stargazers_count"`
}

func mapGitHubResponse(payload githubSearchResponse) []hybridsearch.SearchResult {
	if len(payload.Items) == 0 {
		return []hybridsearch.SearchResult{}
	}

	out := make([]hybridsearch.SearchResult, 0, len(payload.Items))

	for _, item := range payload.Items {
		title := strings.TrimSpace(item.FullName)
		target := strings.TrimSpace(item.HTMLURL)
		if title == "" || target == "" {
			continue
		}

		snippet := strings.TrimSpace(item.Description)
		if item.Language != "" {
			if snippet != "" {
				snippet += " • "
			}
			snippet += "Language: " + item.Language
		}
		if item.StargazersCount > 0 {
			if snippet != "" {
				snippet += " • "
			}
			snippet += fmt.Sprintf("★ %d", item.StargazersCount)
		}

		out = append(out, hybridsearch.SearchResult{
			ID:      strconv.FormatInt(item.ID, 10),
			Title:   title,
			Target:  target,
			Snippet: snippet,
			Source:  "github",
			Kind:    hybridsearch.ResultKindAPI,
			Score:   3,
			Metadata: map[string]string{
				"provider": "github",
				"language": item.Language,
			},
		})
	}

	return out
}
