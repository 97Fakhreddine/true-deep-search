package index

import (
	"context"

	"hybridsearch/internal/search"
)

type Repository interface {
	Index(ctx context.Context, docs []Document) error
	Delete(ctx context.Context, ids []string) error
	Search(ctx context.Context, query string, limit int) ([]search.SearchResult, error)
}
