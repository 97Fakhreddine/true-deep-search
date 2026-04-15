package web

import (
	"context"
	"strings"

	hybridsearch "hybridsearch/pkg/hybridsearch"
)

type Provider struct {
	client Client
	name   string
}

func New(client Client) *Provider {
	return &Provider{
		client: client,
		name:   "web",
	}
}

func (p *Provider) Name() string {
	if p.name == "" {
		return "web"
	}
	return p.name
}

func (p *Provider) Kind() hybridsearch.ResultKind {
	return hybridsearch.ResultKindWeb
}

func (p *Provider) Search(ctx context.Context, req hybridsearch.SearchRequest) ([]hybridsearch.SearchResult, error) {
	query := strings.TrimSpace(req.Query)
	if query == "" {
		return []hybridsearch.SearchResult{}, nil
	}

	if p.client == nil {
		return []hybridsearch.SearchResult{}, nil
	}

	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}

	rawResults, err := p.client.Search(ctx, query, limit)
	if err != nil {
		return nil, err
	}

	return normalizeRemoteResults(rawResults, p.Name()), nil
}
