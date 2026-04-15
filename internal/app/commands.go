package app

import (
	"fmt"
	"hybridsearch/internal/ai"

	tea "github.com/charmbracelet/bubbletea"

	"hybridsearch/internal/search"
)

// 🔎 Normal search
func runSearchCmd(orch *search.Orchestrator, query string) tea.Cmd {
	return func() tea.Msg {
		resp, err := orch.Search(
			nil,
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

func runAICmd(aiService *ai.Service, query string) tea.Cmd {
	return func() tea.Msg {
		answer, err := aiService.Ask(query)

		if err != nil {
			fmt.Println("AI ERROR:", err)
		} else {
			fmt.Println("AI OK")
		}

		return AIFinishedMsg{
			Answer: answer,
			Err:    err,
			Query:  query,
		}
	}
}
