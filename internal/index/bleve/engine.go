package bleve

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	blevev2 "github.com/blevesearch/bleve/v2"

	"hybridsearch/internal/index"
	"hybridsearch/internal/search"
)

type Engine struct {
	index blevev2.Index
	path  string
}

func New(path string) (*Engine, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		path = "./data/index.bleve"
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("create bleve index dir: %w", err)
	}

	var idx blevev2.Index
	var err error

	if _, statErr := os.Stat(path); statErr == nil {
		idx, err = blevev2.Open(path)
		if err != nil {
			return nil, fmt.Errorf("open bleve index: %w", err)
		}
	} else {
		idx, err = blevev2.New(path, newIndexMapping())
		if err != nil {
			return nil, fmt.Errorf("create bleve index: %w", err)
		}
	}

	return &Engine{
		index: idx,
		path:  path,
	}, nil
}

func (e *Engine) Close() error {
	if e == nil || e.index == nil {
		return nil
	}
	return e.index.Close()
}

func (e *Engine) Index(ctx context.Context, docs []index.Document) error {
	if e == nil || e.index == nil {
		return fmt.Errorf("bleve index is nil")
	}

	if len(docs) == 0 {
		return nil
	}

	batch := e.index.NewBatch()

	for _, doc := range docs {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if strings.TrimSpace(doc.ID) == "" {
			continue
		}

		record := storeDocument{
			ID:      doc.ID,
			Path:    doc.Path,
			Title:   doc.Title,
			Content: doc.Content,
		}

		if err := batch.Index(doc.ID, record); err != nil {
			return fmt.Errorf("batch index doc %s: %w", doc.ID, err)
		}
	}

	if err := e.index.Batch(batch); err != nil {
		return fmt.Errorf("commit bleve batch: %w", err)
	}

	return nil
}

func (e *Engine) Delete(ctx context.Context, ids []string) error {
	if e == nil || e.index == nil {
		return fmt.Errorf("bleve index is nil")
	}

	for _, id := range ids {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		id = strings.TrimSpace(id)
		if id == "" {
			continue
		}

		if err := e.index.Delete(id); err != nil {
			return fmt.Errorf("delete doc %s: %w", id, err)
		}
	}

	return nil
}

func (e *Engine) Search(ctx context.Context, query string, limit int) ([]search.SearchResult, error) {
	if e == nil || e.index == nil {
		return nil, fmt.Errorf("bleve index is nil")
	}

	query = strings.TrimSpace(query)
	if query == "" {
		return []search.SearchResult{}, nil
	}

	if limit <= 0 {
		limit = 20
	}

	q := blevev2.NewQueryStringQuery(query)
	req := blevev2.NewSearchRequestOptions(q, limit, 0, false)
	req.Fields = []string{"title", "path", "content"}

	res, err := e.index.SearchInContext(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("bleve search: %w", err)
	}

	return mapSearchResults(res.Hits), nil
}
