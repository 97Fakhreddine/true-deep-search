package hybridsearch

type ResultKind string

const (
	ResultKindWeb   ResultKind = "web"
	ResultKindLocal ResultKind = "local"
	ResultKindAPI   ResultKind = "api"
)

type SearchResult struct {
	ID          string
	Title       string
	Target      string
	Snippet     string
	Source      string
	Kind        ResultKind
	Score       float64
	Category    string
	DisplayHost string
	Timestamp   string
	Metadata    map[string]string
}
type SearchRequest struct {
	Query     string
	Limit     int
	Providers []string
}

type ProviderResultInfo struct {
	Provider    string
	ResultCount int
	DurationMS  int64
	Error       string
}

type SearchResponse struct {
	Query        string
	Results      []SearchResult
	ProviderInfo []ProviderResultInfo
	DurationMS   int64
}
