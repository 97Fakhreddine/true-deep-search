package hybridsearch

import "context"

type Provider interface {
	Name() string
	Kind() ResultKind
	Search(ctx context.Context, req SearchRequest) ([]SearchResult, error)
}

type SearchEngine interface {
	Search(ctx context.Context, req SearchRequest) (SearchResponse, error)
}

type Index interface {
	Index(ctx context.Context, docs []any) error
	Delete(ctx context.Context, ids []string) error
}
