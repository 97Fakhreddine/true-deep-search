<!-- Replace file: README.md -->

# True Deep Search

**True Deep Search** is an open-source terminal-based hybrid deep search engine written in Go.

It is not a simple CLI wrapper around a single API.  
It is designed as a real search system that combines:

- **Meta-search** across multiple remote providers
- **Local search** through indexed content
- **Unified ranking, deduplication, and aggregation**
- **Keyboard-first TUI experience**

The goal is to let you search for almost anything directly from the terminal with a fast, fluid, extensible experience.

---

## Vision

True Deep Search aims to bring a modern search-engine experience into the terminal.

Instead of forcing users to:
- open a browser
- jump between websites
- search every source manually

True Deep Search sends one query to multiple providers, collects results, normalizes them, merges them, removes duplicates, ranks them, and displays them in a unified terminal UI.

It is built to become a **real deep search engine**, not just a command launcher.

---

## Current Features

- Terminal UI built with Bubble Tea
- Search input always visible
- Keyboard-first navigation
- Debounced search
- Non-blocking search execution
- Multi-provider orchestration
- Result normalization into one shared structure
- Deduplication layer
- Ranking layer
- Local search provider
- Web provider
- Wikipedia provider
- GitHub provider
- Stack Exchange provider
- Reddit provider
- YouTube provider
- Open selected result in browser or file

---

## Long-Term Goal

True Deep Search is meant to evolve into a real hybrid search platform with:

- more remote providers
- stronger local indexing
- smarter ranking
- better UI filtering
- provider health visibility
- query intent detection
- semantic and vector-based search in the future

---

## Architecture Overview

The project follows a layered architecture:

### 1. UI Layer
Responsible for:
- rendering
- layout
- keyboard interactions

No business logic should live here.

### 2. App/State Layer
Responsible for:
- state management
- event handling
- command dispatching
- user interaction flow

### 3. Search Orchestrator
Responsible for:
- coordinating providers
- running concurrent searches
- collecting provider outputs
- handling partial failures
- sending results through the pipeline

### 4. Provider Layer
Responsible for:
- remote providers
- local providers
- source-specific normalization

Examples:
- web
- wikipedia
- github
- stackexchange
- reddit
- youtube
- local

### 5. Indexing Layer
Responsible for:
- local content extraction
- indexing
- filesystem watching
- local retrieval

### 6. Aggregation Layer
Responsible for:
- merging all provider outputs into one stream

### 7. Deduplication Layer
Responsible for:
- removing duplicate results
- reducing repeated links/titles

### 8. Ranking Layer
Responsible for:
- scoring results
- ordering results by relevance
- source-aware and intent-aware weighting

### 9. Infrastructure Layer
Responsible for:
- browser opening
- config loading
- HTTP clients
- logging
- debounce helpers

---

## Folder Structure

```text
hybridsearch/
├── cmd/
│   └── hybridsearch/
│       └── main.go
│
├── internal/
│   ├── app/
│   ├── tui/
│   ├── search/
│   ├── provider/
│   ├── index/
│   ├── aggregate/
│   ├── dedupe/
│   ├── rank/
│   ├── infra/
│   └── platform/
│
├── pkg/
│   └── hybridsearch/
│
├── configs/
├── docs/
├── scripts/
├── test/
└── README.md
```

---

## Search Flow

The search flow is:

1. User types a query
2. Input is debounced
3. Search request is sent to the orchestrator
4. Orchestrator fans out the query to multiple providers concurrently
5. Providers return normalized results
6. Results are merged
7. Duplicates are removed
8. Results are ranked
9. Final results are rendered in the TUI

---

## Core Principles

- **Fast startup**
- **Responsive UI**
- **Extensible provider model**
- **Minimal coupling**
- **Clean interfaces**
- **Production-friendly architecture**
- **Contributor-friendly project layout**

---

## Tech Stack

- **Go**
- **Bubble Tea**
- **Bubbles**
- **Lip Gloss**
- **Bleve**
- **fsnotify**
- **BurntSushi/toml**

---

## Installation

### 1. Clone the repository

```bash
git clone https://github.com/YOUR_USERNAME/true-deep-search.git
cd true-deep-search
```

### 2. Install dependencies

```bash
go mod tidy
```

### 3. Run the app

```bash
go run ./cmd/hybridsearch
```

---

## Development Scripts

Run in development mode:

```bash
./scripts/dev.sh
```

Lint the project:

```bash
./scripts/lint.sh
```

Run tests:

```bash
./scripts/test.sh
```

---

## Configuration

Example config file:

```toml
[app]
result_limit = 20
debounce_ms = 250

[providers]
enabled = ["local"]

[providers.web]
timeout_ms = 1200

[index]
path = "./data/index.bleve"
watch = false
roots = [
  "./docs",
  "./notes"
]
```

---

## Current Providers

### General Search
- Web
- Wikipedia

### Developer Search
- GitHub
- Stack Exchange

### Community Search
- Reddit

### Video Search
- YouTube

### Local Search
- Bleve-based local index

---

## Keyboard Controls

- `↑` / `k` → move up
- `↓` / `j` → move down
- `Enter` → open selected result
- `Esc` → close
- `Ctrl+C` → quit

---

## Project Status

This project is actively evolving.

Current stage:
- working TUI
- multi-provider search
- hybrid orchestration
- local + remote architecture
- result ranking and deduplication

Still planned:
- stronger provider coverage
- smarter ranking
- richer local indexing
- better UI filters
- split layout / preview panel
- improved caching
- provider health indicators
- semantic search extensions

---

## Roadmap

### Phase 1
- core TUI
- orchestrator
- provider registry
- ranking
- dedupe
- local indexing foundation

### Phase 2
- real external providers
- wikipedia
- github
- stackexchange
- reddit
- youtube

### Phase 3
- improved ranking
- intent-aware scoring
- better UI and filters
- better provider visibility

### Phase 4
- richer indexing
- smarter previews
- semantic capabilities
- advanced search modes

---

## Why This Project Exists

Search in the terminal is usually fragmented.

You either:
- use a browser
- use a site-specific CLI
- or use a narrow wrapper around one service

True Deep Search exists to explore a better model:
a **real hybrid deep search engine** with a terminal-native UX.

---

## Contributing

Contributions are welcome.

Good contributions include:
- new providers
- ranking improvements
- indexing improvements
- UI/UX improvements
- tests
- docs
- bug fixes

Before opening a PR:

1. Keep changes focused
2. Follow the existing architecture
3. Avoid putting business logic in the TUI layer
4. Add or update tests when relevant
5. Update docs if the behavior changes

---

## Development Guidelines

- Prefer small, focused packages
- Keep interfaces clean
- Normalize provider results before returning them
- Use context for network and indexing operations
- Keep the UI responsive
- Prefer standard library solutions unless extra dependencies are clearly worth it

---

## Potential Future Providers

- Hacker News
- PubMed
- MDN
- Dev.to
- News providers
- Documentation providers
- Archive.org
- domain-specific research providers

---

## License

MIT

---
## Maintainer

This project is actively maintained by the original author.

For major changes, please open an issue first to discuss your ideas.

## Author

**Fakhreddine Messaoudi**

- GitHub: https://github.com/97Fakhreddine

True Deep Search is designed and built as an open-source hybrid deep search engine with a focus on performance, extensibility, and real-world usability.

Made with ambition in Tunisia 🇹🇳