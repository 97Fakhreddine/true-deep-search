package app

import (
	"context"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

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

	lastIssuedQuery string
}

func NewModel() Model {
	input := textinput.New()
	input.Placeholder = "Search anything..."
	input.Focus()
	input.CharLimit = 256
	input.Width = 60

	reg := provider.NewRegistry()

	// Local provider
	bleveEngine, err := bleveindex.New("./data/index.bleve")
	if err == nil {
		_ = reg.Register(localprovider.New(bleveEngine))
	} else {
		_ = reg.Register(localprovider.New(noopLocalIndex{}))
	}

	// General web provider
	_ = reg.Register(webprovider.New(
		webprovider.NewDuckDuckGoLiteClient(1500 * time.Millisecond),
	))

	// Knowledge / API providers
	_ = reg.Register(apiprovider.NewWikipediaProvider(1500 * time.Millisecond))
	_ = reg.Register(apiprovider.NewGitHubProvider(1500*time.Millisecond, ""))
	_ = reg.Register(apiprovider.NewStackExchangeProvider(1500 * time.Millisecond))
	_ = reg.Register(apiprovider.NewRedditProvider(1500 * time.Millisecond))
	_ = reg.Register(apiprovider.NewYouTubeProvider(1500 * time.Millisecond))

	orch := search.NewOrchestrator(reg)
	orch.SetProviderTimeout(1500 * time.Millisecond)
	orch.SetTotalTimeout(3000 * time.Millisecond)

	return Model{
		state: State{},
		input: input,
		orch:  orch,
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
			Snippet: "This is a fallback local result because the real local index is not available.",
			Source:  "local",
			Kind:    search.ResultKindLocal,
			Score:   10,
			Metadata: map[string]string{
				"provider": "local-fallback",
			},
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
