package local

import (
	"context"
	"strings"

	hybridsearch "hybridsearch/pkg/hybridsearch"
)

type LocalIndex interface {
	Search(ctx context.Context, query string, limit int) ([]hybridsearch.SearchResult, error)
}

type Provider struct {
	index LocalIndex
}

func New(index LocalIndex) *Provider {
	return &Provider{
		index: index,
	}
}

func (p *Provider) Name() string {
	return "local"
}

func (p *Provider) Kind() hybridsearch.ResultKind {
	return hybridsearch.ResultKindLocal
}

func (p *Provider) Search(ctx context.Context, req hybridsearch.SearchRequest) ([]hybridsearch.SearchResult, error) {
	query := normalizeQuery(req.Query)
	if query == "" {
		return []hybridsearch.SearchResult{}, nil
	}

	if p.index == nil {
		return []hybridsearch.SearchResult{}, nil
	}

	limit := req.Limit
	if limit <= 0 {
		limit = 20
	}

	results, err := p.index.Search(ctx, query, limit)
	if err != nil {
		return nil, err
	}

	return normalizeResults(results), nil
}

func normalizeQuery(q string) string {
	return strings.TrimSpace(q)
}
