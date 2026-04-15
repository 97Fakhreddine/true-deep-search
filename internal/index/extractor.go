package index

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Extractor interface {
	Supports(path string) bool
	Extract(ctx context.Context, path string) (Document, error)
}

type PlainTextExtractor struct {
	extensions map[string]struct{}
}

func NewPlainTextExtractor(exts []string) *PlainTextExtractor {
	if len(exts) == 0 {
		exts = []string{".txt", ".md", ".log", ".json", ".yaml", ".yml"}
	}

	normalized := make(map[string]struct{}, len(exts))
	for _, ext := range exts {
		ext = strings.ToLower(strings.TrimSpace(ext))
		if ext == "" {
			continue
		}
		if !strings.HasPrefix(ext, ".") {
			ext = "." + ext
		}
		normalized[ext] = struct{}{}
	}

	return &PlainTextExtractor{
		extensions: normalized,
	}
}

func (e *PlainTextExtractor) Supports(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	_, ok := e.extensions[ext]
	return ok
}

func (e *PlainTextExtractor) Extract(ctx context.Context, path string) (Document, error) {
	select {
	case <-ctx.Done():
		return Document{}, ctx.Err()
	default:
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return Document{}, fmt.Errorf("read file %s: %w", path, err)
	}

	title := filepath.Base(path)

	return Document{
		ID:      path,
		Path:    path,
		Title:   title,
		Content: string(data),
		Metadata: map[string]string{
			"extension": strings.ToLower(filepath.Ext(path)),
		},
	}, nil
}
