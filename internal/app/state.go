package app

import "hybridsearch/internal/search"

type State struct {
	Query         string
	Results       []search.SearchResult
	SelectedIndex int
	Loading       bool
	Error         string

	Width  int
	Height int

	HasSearched bool
}
