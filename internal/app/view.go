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

			AIAnswer:  m.state.AIAnswer,
			AILoading: m.state.AILoading,
			AIError:   m.state.AIError,

			Width:       m.state.Width,
			Height:      m.state.Height,
			HasSearched: m.state.HasSearched,

			// viewport-rendered content
			ResultsContent: m.resultsViewport.View(),
			AIContent:      m.aiViewport.View(),
		},
		m.input,
	)
}
