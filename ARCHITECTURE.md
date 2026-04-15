```bash

hybridsearch/
├── cmd/
│   └── hybridsearch/
│       └── main.go
│
├── internal/
│   ├── app/
│   │   ├── model.go
│   │   ├── update.go
│   │   ├── commands.go
│   │   ├── state.go
│   │   └── messages.go
│   │
│   ├── tui/
│   │   ├── view.go
│   │   ├── styles.go
│   │   ├── layout.go
│   │   └── keymap.go
│   │
│   ├── search/
│   │   ├── orchestrator.go
│   │   ├── request.go
│   │   ├── result.go
│   │   ├── pipeline.go
│   │   └── errors.go
│   │
│   ├── provider/
│   │   ├── provider.go
│   │   ├── registry.go
│   │   ├── web/
│   │   │   ├── provider.go
│   │   │   ├── client.go
│   │   │   └── normalize.go
│   │   ├── api/
│   │   │   ├── provider.go
│   │   │   └── normalize.go
│   │   └── local/
│   │       ├── provider.go
│   │       ├── query.go
│   │       └── normalize.go
│   │
│   ├── index/
│   │   ├── indexer.go
│   │   ├── document.go
│   │   ├── watcher.go
│   │   ├── extractor.go
│   │   ├── repository.go
│   │   └── bleve/
│   │       ├── engine.go
│   │       ├── mapping.go
│   │       └── store.go
│   │
│   ├── aggregate/
│   │   └── merger.go
│   │
│   ├── dedupe/
│   │   ├── deduper.go
│   │   └── fingerprint.go
│   │
│   ├── rank/
│   │   ├── ranker.go
│   │   ├── scoring.go
│   │   └── heuristics.go
│   │
│   ├── infra/
│   │   ├── browser/
│   │   │   └── opener.go
│   │   ├── config/
│   │   │   ├── config.go
│   │   │   └── loader.go
│   │   ├── http/
│   │   │   └── client.go
│   │   ├── log/
│   │   │   └── logger.go
│   │   └── debounce/
│   │       └── debounce.go
│   │
│   └── platform/
│       ├── contextutil/
│       │   └── cancelgroup.go
│       └── open/
│           └── target.go
│
├── pkg/
│   └── hybridsearch/
│       ├── types.go
│       └── interfaces.go
│
├── configs/
│   └── config.example.toml
│
├── docs/
│   ├── architecture.md
│   ├── providers.md
│   ├── indexing.md
│   ├── ranking.md
│   └── roadmap.md
│
├── scripts/
│   ├── dev.sh
│   └── lint.sh
│
├── test/
│   ├── integration/
│   ├── fixtures/
│   └── e2e/
│
├── .github/
│   ├── workflows/
│   │   ├── ci.yml
│   │   └── release.yml
│   ├── ISSUE_TEMPLATE/
│   ├── pull_request_template.md
│   └── CODEOWNERS
│
├── go.mod
├── go.sum
├── README.md
├── CONTRIBUTING.md
├── LICENSE
└── Makefile

```