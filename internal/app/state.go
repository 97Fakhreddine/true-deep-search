package app

import "hybridsearch/internal/search"

type State struct {
	// 🔎 Search
	Query         string
	Results       []search.SearchResult
	SelectedIndex int
	Loading       bool
	Error         string

	// 🤖 AI Panel
	AIAnswer  string
	AILoading bool
	AIError   string

	// 📐 UI
	Width  int
	Height int

	HasSearched bool
}
