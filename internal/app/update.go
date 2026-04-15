package app

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"hybridsearch/internal/infra/browser"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.state.Width = msg.Width
		m.state.Height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
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
			m.state.HasSearched = true
			m.state.Error = ""
			m.lastIssuedQuery = query

			return m, runSearchCmd(m.orch, query)
		}

	case SearchFinishedMsg:
		if msg.Query != m.lastIssuedQuery {
			return m, nil
		}

		m.state.Loading = false

		if msg.Err != nil {
			m.state.Error = msg.Err.Error()
			m.state.Results = nil
			return m, nil
		}

		m.state.Results = msg.Response.Results
		m.state.SelectedIndex = 0
		m.state.Error = ""

		return m, nil
	}

	m.input, cmd = m.input.Update(msg)
	return m, cmd
}
