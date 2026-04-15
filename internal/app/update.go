package app

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"hybridsearch/internal/infra/browser"
	"hybridsearch/internal/tui"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.state.Width = msg.Width
		m.state.Height = msg.Height

		// update viewport sizes
		m.syncLayout()

		return m, nil

	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c", "esc", "ctrl+q":
			return m, tea.Quit

		case "up", "k":
			if m.state.SelectedIndex > 0 {
				m.state.SelectedIndex--
			}
			return m, nil

		case "down", "j":
			if m.state.SelectedIndex < len(m.state.Results)-1 {
				m.state.SelectedIndex++
			}
			return m, nil

		case "ctrl+o":
			if len(m.state.Results) > 0 &&
				m.state.SelectedIndex >= 0 &&
				m.state.SelectedIndex < len(m.state.Results) {

				selected := m.state.Results[m.state.SelectedIndex]
				_ = browser.Open(selected.Target)
			}
			return m, nil

		case "enter":
			query := strings.TrimSpace(m.input.Value())
			if query == "" {
				return m, nil
			}

			m.state.Query = query
			m.state.Loading = true
			m.state.Error = ""
			m.state.HasSearched = true
			m.lastIssuedQuery = query

			m.state.AILoading = true
			m.state.AIError = ""
			m.state.AIAnswer = ""

			return m, tea.Batch(
				runSearchCmd(m.orch, query),
				runAICmd(m.ai, query),
			)
		}
	}

	// update input
	m.input, cmd = m.input.Update(msg)

	switch msg := msg.(type) {

	case SearchFinishedMsg:
		if msg.Query != m.lastIssuedQuery {
			return m, nil
		}

		m.state.Loading = false

		if msg.Err != nil {
			m.state.Error = msg.Err.Error()
			m.state.Results = nil
		} else {
			m.state.Results = msg.Response.Results
			m.state.SelectedIndex = 0
			m.state.Error = ""
		}

	case AIFinishedMsg:
		if msg.Query != m.lastIssuedQuery {
			return m, nil
		}

		m.state.AILoading = false

		if msg.Err != nil {
			m.state.AIError = msg.Err.Error()
			m.state.AIAnswer = ""
		} else {
			m.state.AIAnswer = msg.Answer
			m.state.AIError = ""
		}
	}

	// IMPORTANT: update viewport content EVERY FRAME

	// build results text using tui logic
	resultsContent := tui.BuildResultsForViewport(m.state)
	m.setResultsViewportContent(resultsContent)

	// build AI content
	aiContent := tui.BuildAIForViewport(m.state)
	m.setAIViewportContent(aiContent)

	// update scroll behavior
	m.resultsViewport, _ = m.resultsViewport.Update(msg)
	m.aiViewport, _ = m.aiViewport.Update(msg)

	return m, cmd
}
