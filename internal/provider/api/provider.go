package api

import (
	"context"
	"strings"

	hybridsearch "hybridsearch/pkg/hybridsearch"
)

type Client interface {
	Search(ctx context.Context, query string, limit int) ([]APIResult, error)
}

type Provider struct {
	client Client
	name   string
}

func New(name string, client Client) *Provider {
	if strings.TrimSpace(name) == "" {
		name = "api"
	}

	return &Provider{
		client: client,
		name:   name,
	}
}

func (p *Provider) Name() string {
	return p.name
}

func (p *Provider) Kind() hybridsearch.ResultKind {
	return hybridsearch.ResultKindAPI
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
		limit = 20
	}

	rawResults, err := p.client.Search(ctx, query, limit)
	if err != nil {
		return nil, err
	}

	return normalizeAPIResults(rawResults, p.name), nil
}
