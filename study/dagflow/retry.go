package dagflow

import (
	"context"
	"sync/atomic"
	"time"
)

// RetryPolicy configures retry behavior for a node.
type RetryPolicy struct {
	MaxRetries    int                  // max retry attempts (0 = no retry)
	InitialDelay  time.Duration        // delay before first retry
	BackoffFactor float64              // multiplier for each subsequent delay (e.g. 2.0 for exponential)
	Endpoints     []string             // ring buffer of endpoints to rotate through
	IsRetryable   func(err error) bool // optional: determines if an error is retryable
}

// EndpointFunc is a node function that receives the current endpoint address.
type EndpointFunc func(ctx context.Context, deps map[string]any, endpoint string) (any, error)

// ring tracks the current position in the endpoint list.
type ring struct {
	endpoints []string
	pos       atomic.Int64
}

func newRing(endpoints []string) *ring {
	return &ring{endpoints: endpoints}
}

// Next returns the next endpoint in the ring, rotating circularly.
func (r *ring) Next() string {
	idx := r.pos.Add(1) - 1
	return r.endpoints[idx%int64(len(r.endpoints))]
}

// withRetry wraps a NodeFunc with retry logic using the ring buffer.
func withRetry(policy *RetryPolicy, fn NodeFunc) NodeFunc {
	if policy == nil || policy.MaxRetries <= 0 {
		return fn
	}

	return func(ctx context.Context, deps map[string]any) (any, error) {
		var lastErr error
		delay := policy.InitialDelay

		for attempt := 0; attempt <= policy.MaxRetries; attempt++ {
			if attempt > 0 {
				// Wait before retry, respect context cancellation
				select {
				case <-ctx.Done():
					return nil, ctx.Err()
				case <-time.After(delay):
				}
				if policy.BackoffFactor > 0 {
					delay = time.Duration(float64(delay) * policy.BackoffFactor)
				}
			}

			val, err := fn(ctx, deps)
			if err == nil {
				return val, nil
			}
			lastErr = err

			// Check if error is retryable
			if policy.IsRetryable != nil && !policy.IsRetryable(err) {
				return nil, lastErr
			}
		}
		return nil, lastErr
	}
}

// withEndpointRetry wraps an EndpointFunc with retry + ring buffer rotation.
func withEndpointRetry(policy *RetryPolicy, fn EndpointFunc) NodeFunc {
	r := newRing(policy.Endpoints)

	nodeFn := func(ctx context.Context, deps map[string]any) (any, error) {
		endpoint := r.Next()
		return fn(ctx, deps, endpoint)
	}

	return withRetry(policy, nodeFn)
}
