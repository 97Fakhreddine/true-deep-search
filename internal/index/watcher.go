package index

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
)

type Watcher struct {
	indexer *Indexer
	watcher *fsnotify.Watcher
}

func NewWatcher(indexer *Indexer) (*Watcher, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("create fs watcher: %w", err)
	}

	return &Watcher{
		indexer: indexer,
		watcher: w,
	}, nil
}

func (w *Watcher) AddRoots(roots []string) error {
	for _, root := range roots {
		root = strings.TrimSpace(root)
		if root == "" {
			continue
		}

		if err := w.watcher.Add(root); err != nil {
			return fmt.Errorf("watch root %s: %w", root, err)
		}
	}

	return nil
}

func (w *Watcher) Run(ctx context.Context) error {
	if w.indexer == nil {
		return fmt.Errorf("indexer is nil")
	}

	defer w.watcher.Close()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case event, ok := <-w.watcher.Events:
			if !ok {
				return nil
			}

			if event.Op&(fsnotify.Create|fsnotify.Write) != 0 {
				_ = w.indexer.IndexPaths(ctx, []string{event.Name})
			}

			if event.Op&fsnotify.Rename != 0 {
				_ = w.indexer.IndexPaths(ctx, []string{event.Name})
			}

			if event.Op&fsnotify.Remove != 0 {
				_ = w.indexer.repo.Delete(ctx, []string{filepath.Clean(event.Name)})
			}

		case _, ok := <-w.watcher.Errors:
			if !ok {
				return nil
			}
		}
	}
}
