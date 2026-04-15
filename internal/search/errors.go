package search

import "errors"

var (
	ErrEmptyQuery       = errors.New("empty query")
	ErrNoProviders      = errors.New("no providers available")
	ErrProviderNotFound = errors.New("provider not found")
	ErrInvalidSearchReq = errors.New("invalid search request")
)
