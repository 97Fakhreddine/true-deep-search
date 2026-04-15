package provider

import (
	"context"

	hybridsearch "hybridsearch/pkg/hybridsearch"
)

type Provider interface {
	Name() string
	Kind() hybridsearch.ResultKind
	Search(ctx context.Context, req hybridsearch.SearchRequest) ([]hybridsearch.SearchResult, error)
}
