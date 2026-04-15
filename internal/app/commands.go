package app

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"hybridsearch/internal/search"
)

const debounceDelay = 250 * time.Millisecond

func debounceSearchCmd(query string) tea.Cmd {
	return tea.Tick(debounceDelay, func(time.Time) tea.Msg {
		return debounceMsg{Query: query}
	})
}

func runSearchCmd(orch *search.Orchestrator, query string) tea.Cmd {
	return func() tea.Msg {
		resp, err := orch.Search(
			nil, // context will be handled later (MVP simplification)
			search.SearchRequest{
				Query: query,
				Limit: 20,
			},
		)

		return SearchFinishedMsg{
			Response: resp,
			Err:      err,
			Query:    query,
		}
	}
}
