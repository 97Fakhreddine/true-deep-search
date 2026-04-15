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
		"",
		searchBox,
		resultsBox,
		"",
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

	subtitle := "True Deep Search • Made in Tunisia"

	block := lipgloss.JoinVertical(
		lipgloss.Center,
		styles.LogoSubtitle.Render(flag),
		"",
		styles.LogoTitle.Render(strings.TrimRight(title, "\n")),
		styles.LogoSubtitle.Render(subtitle),
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
	boxWidth := clamp(state.Width-12, 60, 140)
	boxHeight := clamp(state.Height-28, 10, 22)

	innerWidth := max(20, boxWidth-4)
	body := buildResultsBody(styles, state, innerWidth, boxHeight)

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
		"%s navigate • %s search • %s open • %s close • %s quit",
		keys.Up+"/"+keys.Down,
		keys.Search,
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

func buildResultsBody(styles Styles, state ViewState, innerWidth int, boxHeight int) string {
	if state.Error != "" {
		return styles.Error.Width(innerWidth).Render("Error: " + state.Error)
	}

	if state.Loading && len(state.Results) == 0 {
		return styles.Status.Width(innerWidth).Render("Searching...")
	}

	if len(state.Results) == 0 {
		if strings.TrimSpace(state.Query) == "" {
			return styles.Empty.Width(innerWidth).Render("Search across web, code, videos and communities...")
		}
		return styles.Empty.Width(innerWidth).Render("No results found.")
	}

	visibleCount := computeVisibleItemCount(boxHeight)
	start, end := computeWindow(len(state.Results), state.SelectedIndex, visibleCount)

	lines := make([]string, 0, end-start+2)

	if start > 0 {
		lines = append(lines, styles.Status.Width(innerWidth).Render("↑ more results above"))
	}

	for i := start; i < end; i++ {
		lines = append(lines, renderResultItem(styles, state.Results[i], i == state.SelectedIndex, innerWidth))
	}

	if end < len(state.Results) {
		lines = append(lines, styles.Status.Width(innerWidth).Render("↓ more results below"))
	}

	return strings.Join(lines, "\n")
}

func renderResultItem(styles Styles, result search.SearchResult, selected bool, width int) string {
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

	title = truncate(title, max(20, width-12))
	snippet = truncate(snippet, max(30, width-12))
	target = truncate(target, max(30, width-22))

	badge := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#000000")).
		Background(lipgloss.Color("#A855F7")).
		Padding(0, 1).
		Render(strings.ToUpper(result.Source))

	meta := badge + " " + target

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		styles.Title.Width(width-8).Render(title),
		styles.Snippet.Width(width-8).Render(snippet),
		lipgloss.NewStyle().Width(width-8).Render(meta),
	)

	cardWidth := max(20, width-2)

	cardStyle := lipgloss.NewStyle().
		Width(cardWidth).
		Border(lipgloss.NormalBorder()).
		Padding(0, 1).
		MarginBottom(1)

	if selected {
		cardStyle = cardStyle.
			BorderForeground(lipgloss.Color("#A855F7"))
	}

	card := cardStyle.Render(content)

	if selected {
		marker := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#A855F7")).
			Bold(true).
			Render("▌ ")
		return marker + card
	}

	return "  " + card
}

func truncate(s string, maxLen int) string {
	s = strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(s, "\n", " "), "\r", " "))
	s = strings.Join(strings.Fields(s), " ")

	if maxLen <= 0 || len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}

func computeVisibleItemCount(boxHeight int) int {
	count := (boxHeight - 2) / 5
	return clamp(count, 2, 6)
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

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
