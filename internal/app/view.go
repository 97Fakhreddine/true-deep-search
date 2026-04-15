package app

import "hybridsearch/internal/tui"

func (m Model) View() string {
	return tui.Render(
		tui.ViewState{
			Query:         m.state.Query,
			Results:       m.state.Results,
			SelectedIndex: m.state.SelectedIndex,
			Loading:       m.state.Loading,
			Error:         m.state.Error,
			Width:         m.state.Width,
			Height:        m.state.Height,
			HasSearched:   m.state.HasSearched,
		},
		m.input,
	)
}
