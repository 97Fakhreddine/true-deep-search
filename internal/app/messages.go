package app

import "hybridsearch/internal/search"

type SearchFinishedMsg struct {
	Response search.SearchResponse
	Err      error
	Query    string
}

type debounceMsg struct {
	Query string
}
