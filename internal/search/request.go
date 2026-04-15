package search

import (
	"strings"
	"time"

	hybridsearch "hybridsearch/pkg/hybridsearch"
)

type SearchRequest = hybridsearch.SearchRequest

type ProviderSelection struct {
	Names []string
}

type SearchOptions struct {
	Limit            int
	ProviderTimeout  time.Duration
	TotalTimeout     time.Duration
	IncludeProviders []string
	ExcludeProviders []string
	RequestID        string
}

func NormalizeRequest(req SearchRequest) SearchRequest {
	req.Query = strings.TrimSpace(req.Query)

	if req.Limit <= 0 {
		req.Limit = 20
	}

	return req
}
