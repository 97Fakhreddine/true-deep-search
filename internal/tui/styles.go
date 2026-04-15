package tui

import "github.com/charmbracelet/lipgloss"

type Styles struct {
	App lipgloss.Style

	HeaderBox    lipgloss.Style
	LogoTitle    lipgloss.Style
	LogoSubtitle lipgloss.Style

	SearchBox lipgloss.Style

	ResultsBox     lipgloss.Style
	ResultItem     lipgloss.Style
	SelectedItem   lipgloss.Style
	SelectedMarker lipgloss.Style

	Title   lipgloss.Style
	Snippet lipgloss.Style
	Meta    lipgloss.Style

	Status lipgloss.Style
	Error  lipgloss.Style
	Empty  lipgloss.Style

	KeysBox lipgloss.Style
	Keys    lipgloss.Style
}

func DefaultStyles() Styles {
	return Styles{
		App: lipgloss.NewStyle().
			Padding(2, 4),

		HeaderBox: lipgloss.NewStyle().
			Padding(1, 0, 1, 0),

		LogoTitle: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#A855F7")),

		LogoSubtitle: lipgloss.NewStyle().
			Faint(true).
			Foreground(lipgloss.Color("#9CA3AF")),

		SearchBox: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(0, 1),

		ResultsBox: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(1, 1),

		ResultItem: lipgloss.NewStyle().
			Padding(0, 0, 1, 0),

		SelectedItem: lipgloss.NewStyle().
			Bold(true),

		SelectedMarker: lipgloss.NewStyle().
			Bold(true),

		Title: lipgloss.NewStyle().
			Bold(true),

		Snippet: lipgloss.NewStyle(),

		Meta: lipgloss.NewStyle().
			Faint(true),

		Status: lipgloss.NewStyle().
			Faint(true),

		Error: lipgloss.NewStyle().
			Bold(true),

		Empty: lipgloss.NewStyle().
			Faint(true).
			Italic(true),

		KeysBox: lipgloss.NewStyle().
			Padding(1, 0, 0, 0),

		Keys: lipgloss.NewStyle().
			Faint(true),
	}
}
