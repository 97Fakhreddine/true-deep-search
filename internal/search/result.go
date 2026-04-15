package search

import hybridsearch "hybridsearch/pkg/hybridsearch"

type ResultKind = hybridsearch.ResultKind

const (
	ResultKindWeb   = hybridsearch.ResultKindWeb
	ResultKindLocal = hybridsearch.ResultKindLocal
	ResultKindAPI   = hybridsearch.ResultKindAPI
)

type SearchResult = hybridsearch.SearchResult
type ProviderResultInfo = hybridsearch.ProviderResultInfo
type SearchResponse = hybridsearch.SearchResponse
