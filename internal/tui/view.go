package tui

import (
	"fmt"
	"reflect"
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

	// 🤖 AI
	AIAnswer  string
	AILoading bool
	AIError   string

	// viewport-rendered content
	ResultsContent string
	AIContent      string

	Width       int
	Height      int
	HasSearched bool
}

func Render(state ViewState, input textinput.Model) string {
	styles := DefaultStyles()
	keys := DefaultKeyMap()

	header := renderHeader(styles, state.Width)
	searchBox := renderSearchBox(styles, input, state.Width)
	mainContent := renderSplitContent(styles, state)
	keysBar := renderKeysBar(styles, keys, state.Width)

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		header,
		"",
		searchBox,
		"",
		mainContent,
		"",
		keysBar,
	)

	if !state.HasSearched && len(state.Results) == 0 && !state.Loading && state.Error == "" && strings.TrimSpace(state.Query) == "" {
		return centerContent(state.Width, state.Height, styles.App.Render(content))
	}

	return styles.App.Render(content)
}

func renderSplitContent(styles Styles, state ViewState) string {
	totalWidth := clamp(state.Width-10, 80, 160)
	leftWidth := int(float64(totalWidth) * 0.6)
	rightWidth := totalWidth - leftWidth - 2

	results := renderResultsBox(styles, state, leftWidth)
	ai := renderAIBox(styles, state, rightWidth)

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		results,
		"  ",
		ai,
	)
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

func renderResultsBox(styles Styles, state ViewState, width int) string {
	boxHeight := clamp(state.Height-28, 10, 22)
	innerWidth := max(10, width-4)

	body := state.ResultsContent

	if strings.TrimSpace(body) == "" {
		body = buildResultsBody(styles, state, innerWidth, boxHeight)
	}

	return styles.ResultsBox.
		Width(width).
		Height(boxHeight).
		Render(body)
}

func renderAIBox(styles Styles, state ViewState, width int) string {
	boxHeight := clamp(state.Height-28, 10, 22)
	innerWidth := max(10, width-4)

	content := state.AIContent

	if strings.TrimSpace(content) == "" {
		content = buildAIBody(state, innerWidth)
	}

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1).
		Width(width).
		Height(boxHeight).
		Render(content)
}

func renderKeysBar(styles Styles, keys KeyMap, width int) string {
	content := fmt.Sprintf(
		"%s navigate • %s search • %s open • %s quit",
		keys.Up+"/"+keys.Down,
		keys.Search,
		keys.Open,
		keys.Quit,
	)

	return lipgloss.PlaceHorizontal(
		clamp(width-8, 50, 150),
		lipgloss.Center,
		styles.KeysBox.Render(styles.Keys.Render(content)),
	)
}

func buildResultsBody(styles Styles, state ViewState, width int, height int) string {
	if state.Error != "" {
		return styles.Error.Width(width).Render("Error: " + state.Error)
	}

	if state.Loading && len(state.Results) == 0 {
		return styles.Status.Width(width).Render("Searching...")
	}

	if len(state.Results) == 0 {
		if strings.TrimSpace(state.Query) == "" {
			return styles.Empty.Width(width).Render("Start typing to search...")
		}
		return styles.Empty.Width(width).Render("No results found.")
	}

	lines := make([]string, 0, len(state.Results)+2)

	visibleCount := computeVisibleItemCount(height)
	start, end := computeWindow(len(state.Results), state.SelectedIndex, visibleCount)

	if start > 0 {
		lines = append(lines, styles.Status.Width(width).Render("↑ more results above"))
	}

	for i := start; i < end; i++ {
		lines = append(lines, renderResultItem(styles, state.Results[i], i == state.SelectedIndex, width))
	}

	if end < len(state.Results) {
		lines = append(lines, styles.Status.Width(width).Render("↓ more results below"))
	}

	return strings.Join(lines, "\n")
}

func buildAIBody(state ViewState, width int) string {
	var content string

	switch {
	case state.AIError != "":
		content = "Error:\n" + state.AIError
	case state.AILoading:
		content = "🤖 Thinking...\n\nContacting Gemini..."
	case strings.TrimSpace(state.AIAnswer) != "":
		content = state.AIAnswer
	default:
		content = "AI response will appear here..."
	}

	return wrapAndTruncateParagraphs(content, width)
}

func renderResultItem(styles Styles, result search.SearchResult, selected bool, width int) string {
	title := strings.TrimSpace(result.Title)
	if title == "" {
		title = "(untitled)"
	}

	snippet := strings.TrimSpace(result.Snippet)
	if snippet == "" {
		snippet = "No description available."
	}

	target := strings.TrimSpace(result.Target)
	if target == "" {
		target = "-"
	}

	maxTextWidth := max(10, width-8)

	title = truncate(title, maxTextWidth)
	snippet = truncate(snippet, maxTextWidth)
	target = truncate(target, max(10, maxTextWidth-12))

	badge := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#000")).
		Background(lipgloss.Color("#A855F7")).
		Padding(0, 1).
		Render(strings.ToUpper(result.Source))

	metaLine := badge + " " + target

	item := lipgloss.JoinVertical(
		lipgloss.Left,
		styles.Title.Width(maxTextWidth).Render(title),
		styles.Snippet.Width(maxTextWidth).Render(snippet),
		lipgloss.NewStyle().Width(maxTextWidth).Render(metaLine),
	)

	card := lipgloss.NewStyle().
		Width(width-2).
		Border(lipgloss.NormalBorder()).
		Padding(0, 1).
		MarginBottom(1)

	if selected {
		card = card.BorderForeground(lipgloss.Color("#A855F7"))
		return styles.SelectedItem.Render("▌ " + card.Render(item))
	}

	return "  " + card.Render(item)
}

func BuildResultsForViewport(v any) string {
	state := viewStateFromAny(v)
	return buildResultsBody(DefaultStyles(), state, 72, 18)
}

func BuildAIForViewport(v any) string {
	state := viewStateFromAny(v)
	return buildAIBody(state, 42)
}

func viewStateFromAny(v any) ViewState {
	switch s := v.(type) {
	case ViewState:
		return s
	case *ViewState:
		if s != nil {
			return *s
		}
		return ViewState{}
	default:
		return reflectToViewState(v)
	}
}

func reflectToViewState(v any) ViewState {
	rv := reflect.ValueOf(v)
	if !rv.IsValid() {
		return ViewState{}
	}
	if rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			return ViewState{}
		}
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return ViewState{}
	}

	var out ViewState

	if f := rv.FieldByName("Query"); f.IsValid() && f.Kind() == reflect.String {
		out.Query = f.String()
	}
	if f := rv.FieldByName("SelectedIndex"); f.IsValid() && f.Kind() == reflect.Int {
		out.SelectedIndex = int(f.Int())
	}
	if f := rv.FieldByName("Loading"); f.IsValid() && f.Kind() == reflect.Bool {
		out.Loading = f.Bool()
	}
	if f := rv.FieldByName("Error"); f.IsValid() && f.Kind() == reflect.String {
		out.Error = f.String()
	}
	if f := rv.FieldByName("AIAnswer"); f.IsValid() && f.Kind() == reflect.String {
		out.AIAnswer = f.String()
	}
	if f := rv.FieldByName("AILoading"); f.IsValid() && f.Kind() == reflect.Bool {
		out.AILoading = f.Bool()
	}
	if f := rv.FieldByName("AIError"); f.IsValid() && f.Kind() == reflect.String {
		out.AIError = f.String()
	}
	if f := rv.FieldByName("Width"); f.IsValid() && f.Kind() == reflect.Int {
		out.Width = int(f.Int())
	}
	if f := rv.FieldByName("Height"); f.IsValid() && f.Kind() == reflect.Int {
		out.Height = int(f.Int())
	}
	if f := rv.FieldByName("HasSearched"); f.IsValid() && f.Kind() == reflect.Bool {
		out.HasSearched = f.Bool()
	}
	if f := rv.FieldByName("Results"); f.IsValid() && f.CanInterface() {
		if results, ok := f.Interface().([]search.SearchResult); ok {
			out.Results = results
		}
	}

	return out
}

func wrapAndTruncateParagraphs(s string, width int) string {
	if width <= 0 {
		return s
	}

	paragraphs := strings.Split(strings.ReplaceAll(s, "\r\n", "\n"), "\n")
	out := make([]string, 0, len(paragraphs))

	for _, p := range paragraphs {
		p = strings.TrimSpace(p)
		if p == "" {
			out = append(out, "")
			continue
		}

		words := strings.Fields(p)
		if len(words) == 0 {
			out = append(out, "")
			continue
		}

		line := words[0]
		for _, w := range words[1:] {
			if len(line)+1+len(w) > width {
				out = append(out, line)
				line = w
				continue
			}
			line += " " + w
		}
		out = append(out, line)
	}

	return strings.Join(out, "\n")
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

func computeVisibleItemCount(height int) int {
	count := (height - 2) / 5
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
