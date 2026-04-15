package search

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"

	"hybridsearch/internal/aggregate"
	"hybridsearch/internal/dedupe"
	"hybridsearch/internal/provider"
	"hybridsearch/internal/rank"
)

const (
	DefaultProviderTimeout = 1500 * time.Millisecond
	DefaultTotalTimeout    = 2500 * time.Millisecond
)

type Orchestrator struct {
	registry        *provider.Registry
	aggregator      Aggregator
	deduper         Deduper
	ranker          Ranker
	providerTimeout time.Duration
	totalTimeout    time.Duration
}

func NewOrchestrator(registry *provider.Registry) *Orchestrator {
	return &Orchestrator{
		registry:        registry,
		aggregator:      aggregate.NewMerger(),
		deduper:         dedupe.New(),
		ranker:          rank.New(),
		providerTimeout: DefaultProviderTimeout,
		totalTimeout:    DefaultTotalTimeout,
	}
}

func (o *Orchestrator) SetProviderTimeout(timeout time.Duration) {
	if timeout > 0 {
		o.providerTimeout = timeout
	}
}

func (o *Orchestrator) SetTotalTimeout(timeout time.Duration) {
	if timeout > 0 {
		o.totalTimeout = timeout
	}
}

func (o *Orchestrator) Search(ctx context.Context, req SearchRequest) (SearchResponse, error) {
	start := time.Now()

	if ctx == nil {
		ctx = context.Background()
	}

	req = NormalizeRequest(req)

	if strings.TrimSpace(req.Query) == "" {
		return emptyResponse(req, start), ErrEmptyQuery
	}

	if o.totalTimeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, o.totalTimeout)
		defer cancel()
	}

	providers := o.resolveProviders(req.Providers)
	if len(providers) == 0 {
		return emptyResponse(req, start), ErrNoProviders
	}

	type providerOutput struct {
		name    string
		results []SearchResult
		err     error
		dur     time.Duration
	}

	outCh := make(chan providerOutput, len(providers))
	var wg sync.WaitGroup

	for _, p := range providers {
		p := p
		wg.Add(1)

		go func() {
			defer wg.Done()

			pCtx := ctx
			cancel := func() {}

			if o.providerTimeout > 0 {
				pCtx, cancel = context.WithTimeout(ctx, o.providerTimeout)
			}
			defer cancel()

			pStart := time.Now()
			results, err := p.Search(pCtx, req)

			if results == nil {
				results = []SearchResult{}
			}

			select {
			case outCh <- providerOutput{
				name:    p.Name(),
				results: results,
				err:     err,
				dur:     time.Since(pStart),
			}:
			case <-ctx.Done():
				return
			}
		}()
	}

	go func() {
		wg.Wait()
		close(outCh)
	}()

	mergedInputs := make(map[string][]SearchResult, len(providers))
	providerInfo := make([]ProviderResultInfo, 0, len(providers))

	for item := range outCh {
		mergedInputs[item.name] = item.results

		info := ProviderResultInfo{
			Provider:    item.name,
			ResultCount: len(item.results),
			DurationMS:  item.dur.Milliseconds(),
		}

		if item.err != nil {
			info.Error = item.err.Error()
		}

		providerInfo = append(providerInfo, info)
	}

	results := o.aggregator.Merge(mergedInputs)
	results = o.deduper.Dedupe(results)
	results = o.ranker.Rank(req.Query, results)

	if req.Limit > 0 && len(results) > req.Limit {
		results = results[:req.Limit]
	}

	resp := SearchResponse{
		Query:        req.Query,
		Results:      results,
		ProviderInfo: providerInfo,
		DurationMS:   time.Since(start).Milliseconds(),
	}

	if len(results) > 0 {
		return resp, nil
	}

	if ctx.Err() != nil && !errors.Is(ctx.Err(), context.Canceled) {
		return resp, ctx.Err()
	}

	return resp, nil
}

func (o *Orchestrator) resolveProviders(names []string) []provider.Provider {
	if len(names) == 0 {
		return o.registry.List()
	}

	resolved := make([]provider.Provider, 0, len(names))
	seen := make(map[string]struct{}, len(names))

	for _, name := range names {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}

		if _, ok := seen[name]; ok {
			continue
		}

		p, ok := o.registry.Get(name)
		if !ok {
			continue
		}

		seen[name] = struct{}{}
		resolved = append(resolved, p)
	}

	return resolved
}

func emptyResponse(req SearchRequest, start time.Time) SearchResponse {
	return SearchResponse{
		Query:        req.Query,
		Results:      []SearchResult{},
		ProviderInfo: []ProviderResultInfo{},
		DurationMS:   time.Since(start).Milliseconds(),
	}
}
