package contextutil

import (
	"context"
	"sync"
)

type CancelGroup struct {
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func New(parent context.Context) *CancelGroup {
	if parent == nil {
		parent = context.Background()
	}

	ctx, cancel := context.WithCancel(parent)

	return &CancelGroup{
		ctx:    ctx,
		cancel: cancel,
	}
}

func (g *CancelGroup) Go(fn func(ctx context.Context)) {
	g.wg.Add(1)

	go func() {
		defer g.wg.Done()
		fn(g.ctx)
	}()
}

func (g *CancelGroup) Cancel() {
	g.cancel()
}

func (g *CancelGroup) Wait() {
	g.wg.Wait()
}

func (g *CancelGroup) Context() context.Context {
	return g.ctx
}
