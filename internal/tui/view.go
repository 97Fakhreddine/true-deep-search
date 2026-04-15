package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"

	"hybridsearch/internal/search"
)

type ViewState struct {
	Query         string
	Results       []search.SearchResult
	SelectedIndex int
	Loading       bool
	Error         string
	Width         int
	Height        int
	HasSearched   bool
}

func Render(state ViewState, input textinput.Model) string {
	styles := DefaultStyles()
	keys := DefaultKeyMap()

	header := renderHeader(styles, state.Width)
	searchBox := renderSearchBox(styles, input, state.Width)
	resultsBox := renderResultsBox(styles, state)
	keysBar := renderKeysBar(styles, keys, state.Width)

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		header,
		searchBox,
		"",
		resultsBox,
		keysBar,
	)

	if !state.HasSearched && len(state.Results) == 0 && !state.Loading && state.Error == "" && strings.TrimSpace(state.Query) == "" {
		return centerContent(state.Width, state.Height, styles.App.Render(content))
	}

	return styles.App.Render(content)
}

func renderHeader(styles Styles, width int) string {
	flag := "🇹🇳"

	title := `
████████╗██████╗ ██╗   ██╗███████╗    ██████╗ ███████╗███████╗██████╗ 
╚══██╔══╝██╔══██╗██║   ██║██╔════╝    ██╔══██╗██╔════╝██╔════╝██╔══██╗
   ██║   ██████╔╝██║   ██║█████╗      ██║  ██║█████╗  █████╗  ██████╔╝
   ██║   ██╔══██╗██║   ██║██╔══╝      ██║  ██║██╔══╝  ██╔══╝  ██╔═══╝ 
   ██║   ██║  ██║╚██████╔╝███████╗    ██████╔╝███████╗███████╗██║     
   ╚═╝   ╚═╝  ╚═╝ ╚═════╝ ╚══════╝    ╚═════╝ ╚══════╝╚══════╝╚═╝     
`

	subtitle := "True Deep Search"

	block := lipgloss.JoinVertical(
		lipgloss.Center,

		// flag
		styles.LogoSubtitle.Render(flag),

		// spacing
		"",

		// logo
		styles.LogoTitle.Render(strings.TrimRight(title, "\n")),

		// spacing
		"",
		styles.LogoSubtitle.Render(subtitle+" • Made in Tunisia"),
	)

	return lipgloss.PlaceHorizontal(
		clamp(width-8, 60, 140),
		lipgloss.Center,
		styles.HeaderBox.Render(block),
	)
}

func renderSearchBox(styles Styles, input textinput.Model, width int) string {
	boxWidth := clamp(width-12, 50, 120)
	box := styles.SearchBox.Width(boxWidth).Render(input.View())

	return lipgloss.PlaceHorizontal(
		clamp(width-8, 50, 140),
		lipgloss.Center,
		box,
	)
}

func renderResultsBox(styles Styles, state ViewState) string {
	body := buildResultsBody(styles, state)

	boxWidth := clamp(state.Width-12, 60, 140)
	boxHeight := clamp(state.Height-22, 10, 30)

	box := styles.ResultsBox.
		Width(boxWidth).
		Height(boxHeight).
		Render(body)

	return lipgloss.PlaceHorizontal(
		clamp(state.Width-8, 60, 150),
		lipgloss.Center,
		box,
	)
}

func renderKeysBar(styles Styles, keys KeyMap, width int) string {
	content := fmt.Sprintf(
		"%s navigate • %s open • %s close • %s quit",
		keys.Up+"/"+keys.Down,
		keys.Open,
		keys.Close,
		keys.Quit,
	)

	bar := styles.Keys.Render(content)

	return lipgloss.PlaceHorizontal(
		clamp(width-8, 50, 150),
		lipgloss.Center,
		styles.KeysBox.Render(bar),
	)
}

func buildResultsBody(styles Styles, state ViewState) string {
	if state.Error != "" {
		return styles.Error.Render("Error: " + state.Error)
	}

	if state.Loading && len(state.Results) == 0 {
		return styles.Status.Render("Searching...")
	}

	if len(state.Results) == 0 {
		if strings.TrimSpace(state.Query) == "" {
			return styles.Empty.Render("Start typing to search...")
		}
		return styles.Empty.Render("No results found.")
	}

	visibleCount := computeVisibleItemCount(state.Height)
	start, end := computeWindow(len(state.Results), state.SelectedIndex, visibleCount)

	lines := make([]string, 0, end-start+2)

	if start > 0 {
		lines = append(lines, styles.Status.Render("↑ more results above"))
	}

	for i := start; i < end; i++ {
		lines = append(lines, renderResultItem(styles, state.Results[i], i == state.SelectedIndex))
	}

	if end < len(state.Results) {
		lines = append(lines, styles.Status.Render("↓ more results below"))
	}

	return strings.Join(lines, "\n")
}

func renderResultItem(styles Styles, result search.SearchResult, selected bool) string {
	title := strings.TrimSpace(result.Title)
	if title == "" {
		title = "(untitled)"
	}

	target := strings.TrimSpace(result.Target)
	if target == "" {
		target = "-"
	}

	snippet := strings.TrimSpace(result.Snippet)
	if snippet == "" {
		snippet = "No description available."
	}

	meta := fmt.Sprintf("[%s] %s", strings.ToUpper(result.Source), target)

	item := lipgloss.JoinVertical(
		lipgloss.Left,
		styles.Title.Render(title),
		styles.Snippet.Render(snippet),
		styles.Meta.Render(meta),
	)

	item = styles.ResultItem.Render(item)

	if selected {
		return styles.SelectedItem.Render(
			styles.SelectedMarker.Render("› ") + indentLines(item, "  "),
		)
	}

	return "  " + indentLines(item, "  ")
}

func indentLines(s, prefix string) string {
	parts := strings.Split(s, "\n")
	for i := range parts {
		if i == 0 {
			continue
		}
		parts[i] = prefix + parts[i]
	}
	return strings.Join(parts, "\n")
}

func computeVisibleItemCount(height int) int {
	count := (height - 24) / 4
	return clamp(count, 3, 8)
}

func computeWindow(total, selected, visible int) (int, int) {
	if total <= 0 {
		return 0, 0
	}

	if visible >= total {
		return 0, total
	}

	if selected < 0 {
		selected = 0
	}
	if selected >= total {
		selected = total - 1
	}

	half := visible / 2
	start := selected - half
	if start < 0 {
		start = 0
	}

	end := start + visible
	if end > total {
		end = total
		start = end - visible
		if start < 0 {
			start = 0
		}
	}

	return start, end
}
