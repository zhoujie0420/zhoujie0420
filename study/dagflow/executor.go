package dagflow

import (
	"context"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

// Option configures the executor.
type Option func(*executorConfig)

type executorConfig struct {
	workers int
}

// WithWorkers sets the max number of concurrent goroutines.
// Defaults to the number of nodes if not set.
func WithWorkers(n int) Option {
	return func(c *executorConfig) {
		c.workers = n
	}
}

// Result holds the output of all executed nodes.
type Result struct {
	mu     sync.RWMutex
	data   map[string]any
	errors map[string]error // errors from non-critical nodes
}

// Get returns the result of a specific node.
func (r *Result) Get(name string) (any, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	v, ok := r.data[name]
	return v, ok
}

// MustGet returns the result of a specific node, panics if not found.
func (r *Result) MustGet(name string) any {
	v, ok := r.Get(name)
	if !ok {
		panic(fmt.Sprintf("dagflow: result for node %q not found", name))
	}
	return v
}

// Errors returns all non-critical node errors.
func (r *Result) Errors() map[string]error {
	r.mu.RLock()
	defer r.mu.RUnlock()
	cp := make(map[string]error, len(r.errors))
	for k, v := range r.errors {
		cp[k] = v
	}
	return cp
}

func (r *Result) set(name string, value any) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[name] = value
}

func (r *Result) setError(name string, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.errors[name] = err
}

// Execute runs all nodes in the DAG respecting dependency order.
// It uses a goroutine pool for concurrency control and errgroup for error propagation.
// If any node fails, the context is cancelled and remaining nodes are skipped.
func Execute(ctx context.Context, dag *DAG, opts ...Option) (*Result, error) {
	if err := dag.Validate(); err != nil {
		return nil, err
	}

	cfg := &executorConfig{}
	for _, opt := range opts {
		opt(cfg)
	}

	dag.mu.RLock()
	defer dag.mu.RUnlock()

	nodeCount := len(dag.nodes)
	if nodeCount == 0 {
		return &Result{data: make(map[string]any), errors: make(map[string]error)}, nil
	}

	if cfg.workers <= 0 {
		cfg.workers = nodeCount
	}

	// Build in-degree map and downstream adjacency list
	inDegree := make(map[string]int, nodeCount)
	downstream := make(map[string][]string, nodeCount)

	for name := range dag.nodes {
		inDegree[name] = 0
	}
	for _, node := range dag.nodes {
		for _, dep := range node.deps {
			downstream[dep] = append(downstream[dep], node.name)
			inDegree[node.name]++
		}
	}

	result := &Result{
		data:   make(map[string]any, nodeCount),
		errors: make(map[string]error),
	}
	pool := NewPool(cfg.workers)

	g, gCtx := errgroup.WithContext(ctx)

	// Each node gets a "done" channel that closes when the node finishes successfully
	doneChans := make(map[string]chan struct{}, nodeCount)
	for name := range dag.nodes {
		doneChans[name] = make(chan struct{})
	}

	for _, node := range dag.nodes {
		node := node
		g.Go(func() error {
			// Wait for all dependencies to complete
			for _, dep := range node.deps {
				select {
				case <-gCtx.Done():
					return gCtx.Err()
				case <-doneChans[dep]:
				}
			}

			// Skip if context already cancelled
			if gCtx.Err() != nil {
				return gCtx.Err()
			}

			// Execute via pool to limit concurrency
			errCh := make(chan error, 1)
			pool.Submit(func() {
				if node.fn == nil {
					errCh <- nil
					return
				}

				// Build deps map
				deps := make(map[string]any, len(node.deps))
				for _, dep := range node.deps {
					v, _ := result.Get(dep)
					deps[dep] = v
				}

				// Wrap with retry if configured
				fn := node.fn
				if node.retry != nil && len(node.retry.Endpoints) == 0 {
					fn = withRetry(node.retry, node.fn)
				}

				val, err := fn(gCtx, deps)
				if err != nil {
					errCh <- fmt.Errorf("dagflow: node %q failed: %w", node.name, err)
					return
				}
				result.set(node.name, val)
				errCh <- nil
			})

			select {
			case <-gCtx.Done():
				return gCtx.Err()
			case err := <-errCh:
				if err != nil {
					if node.critical {
						return err // cancel entire DAG
					}
					// non-critical: record error, continue execution
					result.setError(node.name, err)
				}
			}

			// Signal downstream nodes (even if non-critical failed, let downstream decide)
			close(doneChans[node.name])
			return nil
		})
	}

	err := g.Wait()
	pool.Shutdown()

	if err != nil {
		return nil, err
	}
	return result, nil
}
