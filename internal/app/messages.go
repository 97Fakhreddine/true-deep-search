package app

import "hybridsearch/internal/search"

// 🔎 Search result message
type SearchFinishedMsg struct {
	Response search.SearchResponse
	Err      error
	Query    string
}

// 🤖 AI result message
type AIFinishedMsg struct {
	Answer string
	Err    error
	Query  string
}

type debounceMsg struct {
	Query string
}
