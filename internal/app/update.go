package app

import (
	tea "github.com/charmbracelet/bubbletea"

	"hybridsearch/internal/infra/browser"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.state.Width = msg.Width
		m.state.Height = msg.Height

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

		case "enter":
			if len(m.state.Results) > 0 && m.state.SelectedIndex >= 0 && m.state.SelectedIndex < len(m.state.Results) {
				selected := m.state.Results[m.state.SelectedIndex]
				_ = browser.Open(selected.Target)
			}
			return m, nil
		}
	}

	m.input, cmd = m.input.Update(msg)
	m.state.Query = m.input.Value()

	switch msg := msg.(type) {
	case debounceMsg:
		if msg.Query == "" || msg.Query != m.input.Value() {
			return m, nil
		}

		m.state.Loading = true
		m.state.HasSearched = true
		m.lastIssuedQuery = msg.Query

		return m, runSearchCmd(m.orch, msg.Query)

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

	return m, tea.Batch(
		cmd,
		debounceSearchCmd(m.input.Value()),
	)
}
