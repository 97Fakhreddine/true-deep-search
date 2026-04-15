package app

import (
	"context"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"

	"hybridsearch/internal/ai"
	bleveindex "hybridsearch/internal/index/bleve"
	"hybridsearch/internal/provider"
	apiprovider "hybridsearch/internal/provider/api"
	localprovider "hybridsearch/internal/provider/local"
	webprovider "hybridsearch/internal/provider/web"
	"hybridsearch/internal/search"
)

type Model struct {
	state State
	input textinput.Model

	orch *search.Orchestrator
	ai   *ai.Service

	resultsViewport viewport.Model
	aiViewport      viewport.Model

	lastIssuedQuery string
}

func NewModel() Model {
	input := textinput.New()
	input.Placeholder = "Search anything..."
	input.Focus()
	input.CharLimit = 256
	input.Width = 60

	reg := provider.NewRegistry()

	bleveEngine, err := bleveindex.New("./data/index.bleve")
	if err == nil {
		_ = reg.Register(localprovider.New(bleveEngine))
	} else {
		_ = reg.Register(localprovider.New(noopLocalIndex{}))
	}

	_ = reg.Register(webprovider.New(
		webprovider.NewDuckDuckGoLiteClient(1500 * time.Millisecond),
	))

	_ = reg.Register(apiprovider.NewWikipediaProvider(1500 * time.Millisecond))
	_ = reg.Register(apiprovider.NewGitHubProvider(1500*time.Millisecond, ""))
	_ = reg.Register(apiprovider.NewStackExchangeProvider(1500 * time.Millisecond))
	_ = reg.Register(apiprovider.NewRedditProvider(1500 * time.Millisecond))
	_ = reg.Register(apiprovider.NewYouTubeProvider(1500 * time.Millisecond))

	orch := search.NewOrchestrator(reg)
	orch.SetProviderTimeout(1500 * time.Millisecond)
	orch.SetTotalTimeout(3000 * time.Millisecond)

	resultsVP := viewport.New(0, 0)
	resultsVP.SetContent("")

	aiVP := viewport.New(0, 0)
	aiVP.SetContent("")

	return Model{
		state: State{},
		input: input,
		orch:  orch,
		ai:    ai.NewService(),

		resultsViewport: resultsVP,
		aiViewport:      aiVP,
	}
}

type noopLocalIndex struct{}

func (noopLocalIndex) Search(ctx context.Context, query string, limit int) ([]search.SearchResult, error) {
	_ = ctx
	_ = limit

	if query == "" {
		return []search.SearchResult{}, nil
	}

	return []search.SearchResult{
		{
			ID:      "local-1",
			Title:   "Local result: " + query,
			Target:  "/tmp/file.txt",
			Snippet: "Fallback local result (no index available)",
			Source:  "local",
			Kind:    search.ResultKindLocal,
			Score:   10,
		},
	}, nil
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) State() State {
	return m.state
}

func (m Model) Input() textinput.Model {
	return m.input
}

func (m Model) ResultsViewport() viewport.Model {
	return m.resultsViewport
}

func (m Model) AIViewport() viewport.Model {
	return m.aiViewport
}

func (m *Model) syncLayout() {
	if m.state.Width <= 0 || m.state.Height <= 0 {
		return
	}

	totalWidth := clampInt(m.state.Width-10, 80, 160)
	leftWidth := int(float64(totalWidth) * 0.6)
	rightWidth := totalWidth - leftWidth - 2

	panelHeight := clampInt(m.state.Height-28, 10, 22)

	resultsInnerWidth := maxInt(10, leftWidth-4)
	aiInnerWidth := maxInt(10, rightWidth-4)
	innerHeight := maxInt(3, panelHeight-2)

	m.resultsViewport.Width = resultsInnerWidth
	m.resultsViewport.Height = innerHeight

	m.aiViewport.Width = aiInnerWidth
	m.aiViewport.Height = innerHeight
}

func (m *Model) setResultsViewportContent(content string) {
	content = strings.TrimRight(content, "\n")
	m.resultsViewport.SetContent(content)
}

func (m *Model) setAIViewportContent(content string) {
	content = strings.TrimRight(content, "\n")
	m.aiViewport.SetContent(content)
}

func clampInt(v, min, max int) int {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
