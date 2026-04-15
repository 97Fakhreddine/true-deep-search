package index

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type Indexer struct {
	repo       Repository
	extractors []Extractor
}

func NewIndexer(repo Repository, extractors ...Extractor) *Indexer {
	return &Indexer{
		repo:       repo,
		extractors: extractors,
	}
}

func (i *Indexer) IndexPaths(ctx context.Context, roots []string) error {
	if i.repo == nil {
		return fmt.Errorf("repository is nil")
	}

	if len(roots) == 0 {
		return nil
	}

	var docs []Document

	for _, root := range roots {
		root = strings.TrimSpace(root)
		if root == "" {
			continue
		}

		info, err := os.Stat(root)
		if err != nil {
			return fmt.Errorf("stat root %s: %w", root, err)
		}

		if info.IsDir() {
			if err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}

				select {
				case <-ctx.Done():
					return ctx.Err()
				default:
				}

				if d.IsDir() {
					return nil
				}

				doc, ok, err := i.extractDocument(ctx, path)
				if err != nil {
					return err
				}
				if !ok {
					return nil
				}

				docs = append(docs, doc)
				return nil
			}); err != nil {
				return err
			}

			continue
		}

		doc, ok, err := i.extractDocument(ctx, root)
		if err != nil {
			return err
		}
		if ok {
			docs = append(docs, doc)
		}
	}

	if len(docs) == 0 {
		return nil
	}

	return i.repo.Index(ctx, docs)
}

func (i *Indexer) extractDocument(ctx context.Context, path string) (Document, bool, error) {
	for _, extractor := range i.extractors {
		if extractor == nil {
			continue
		}

		if !extractor.Supports(path) {
			continue
		}

		doc, err := extractor.Extract(ctx, path)
		if err != nil {
			return Document{}, false, err
		}

		if strings.TrimSpace(doc.ID) == "" {
			doc.ID = path
		}
		if strings.TrimSpace(doc.Path) == "" {
			doc.Path = path
		}
		if strings.TrimSpace(doc.Title) == "" {
			doc.Title = filepath.Base(path)
		}
		if doc.Metadata == nil {
			doc.Metadata = map[string]string{}
		}

		return doc, true, nil
	}

	return Document{}, false, nil
}
