package dagflow

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestBasicDAG(t *testing.T) {
	dag := NewDAG()

	dag.AddNode(NewNode("A", func(ctx context.Context, deps map[string]interface{}) (interface{}, error) {
		return 1, nil
	}))

	dag.AddNode(NewNode("B", func(ctx context.Context, deps map[string]interface{}) (interface{}, error) {
		return 2, nil
	}))

	dag.AddNode(NewNode("C", func(ctx context.Context, deps map[string]interface{}) (interface{}, error) {
		a := deps["A"].(int)
		b := deps["B"].(int)
		return a + b, nil
	}).DependsOn("A", "B"))

	result, err := Execute(context.Background(), dag)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	val := result.MustGet("C").(int)
	if val != 3 {
		t.Fatalf("expected C=3, got %d", val)
	}
}

func TestDependencyResultPassing(t *testing.T) {
	dag := NewDAG()

	dag.AddNode(NewNode("fetch_user", func(ctx context.Context, deps map[string]interface{}) (interface{}, error) {
		return map[string]string{"name": "Zhou Jie", "id": "123"}, nil
	}))

	dag.AddNode(NewNode("fetch_orders", func(ctx context.Context, deps map[string]interface{}) (interface{}, error) {
		user := deps["fetch_user"].(map[string]string)
		return fmt.Sprintf("orders for %s", user["id"]), nil
	}).DependsOn("fetch_user"))

	dag.AddNode(NewNode("fetch_balance", func(ctx context.Context, deps map[string]interface{}) (interface{}, error) {
		user := deps["fetch_user"].(map[string]string)
		return fmt.Sprintf("balance for %s", user["id"]), nil
	}).DependsOn("fetch_user"))

	result, err := Execute(context.Background(), dag)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	orders := result.MustGet("fetch_orders").(string)
	if orders != "orders for 123" {
		t.Fatalf("unexpected orders: %s", orders)
	}

	balance := result.MustGet("fetch_balance").(string)
	if balance != "balance for 123" {
		t.Fatalf("unexpected balance: %s", balance)
	}
}

func TestCycleDetection(t *testing.T) {
	dag := NewDAG()

	dag.AddNode(NewNode("A", nil).DependsOn("B"))
	dag.AddNode(NewNode("B", nil).DependsOn("C"))
	dag.AddNode(NewNode("C", nil).DependsOn("A"))

	_, err := Execute(context.Background(), dag)
	if err == nil {
		t.Fatal("expected cycle error")
	}
}

func TestUnknownDependency(t *testing.T) {
	dag := NewDAG()
	dag.AddNode(NewNode("A", nil).DependsOn("nonexistent"))

	_, err := Execute(context.Background(), dag)
	if err == nil {
		t.Fatal("expected unknown dependency error")
	}
}

func TestNodeFailureCancelsOthers(t *testing.T) {
	dag := NewDAG()

	dag.AddNode(NewNode("fast", func(ctx context.Context, deps map[string]interface{}) (interface{}, error) {
		return nil, errors.New("fast failed")
	}))

	var slowRan atomic.Bool
	dag.AddNode(NewNode("slow", func(ctx context.Context, deps map[string]interface{}) (interface{}, error) {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(2 * time.Second):
			slowRan.Store(true)
			return "done", nil
		}
	}))

	dag.AddNode(NewNode("downstream", func(ctx context.Context, deps map[string]interface{}) (interface{}, error) {
		return nil, nil
	}).DependsOn("fast", "slow"))

	_, err := Execute(context.Background(), dag)
	if err == nil {
		t.Fatal("expected error from fast node")
	}
}

func TestContextTimeout(t *testing.T) {
	dag := NewDAG()

	dag.AddNode(NewNode("slow", func(ctx context.Context, deps map[string]interface{}) (interface{}, error) {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(5 * time.Second):
			return "done", nil
		}
	}))

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	_, err := Execute(ctx, dag)
	if err == nil {
		t.Fatal("expected timeout error")
	}
}

func TestWorkerPoolLimit(t *testing.T) {
	dag := NewDAG()

	var maxConcurrent atomic.Int32
	var current atomic.Int32

	for i := 0; i < 20; i++ {
		name := fmt.Sprintf("node_%d", i)
		dag.AddNode(NewNode(name, func(ctx context.Context, deps map[string]interface{}) (interface{}, error) {
			c := current.Add(1)
			// Track max concurrency
			for {
				old := maxConcurrent.Load()
				if c <= old || maxConcurrent.CompareAndSwap(old, c) {
					break
				}
			}
			time.Sleep(50 * time.Millisecond)
			current.Add(-1)
			return nil, nil
		}))
	}

	_, err := Execute(context.Background(), dag, WithWorkers(3))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	max := maxConcurrent.Load()
	if max > 3 {
		t.Fatalf("expected max concurrency <= 3, got %d", max)
	}
}

func TestDuplicateNode(t *testing.T) {
	dag := NewDAG()
	dag.AddNode(NewNode("A", nil))
	err := dag.AddNode(NewNode("A", nil))
	if err == nil {
		t.Fatal("expected duplicate node error")
	}
}

func TestEmptyDAG(t *testing.T) {
	dag := NewDAG()
	result, err := Execute(context.Background(), dag)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
}

func TestCriticalNodeFailureCancelsDAG(t *testing.T) {
	dag := NewDAG()

	dag.AddNode(NewNode("critical_node", func(ctx context.Context, deps map[string]any) (any, error) {
		return nil, errors.New("critical failure")
	})) // critical by default

	var downstreamRan atomic.Bool
	dag.AddNode(NewNode("downstream", func(ctx context.Context, deps map[string]any) (any, error) {
		downstreamRan.Store(true)
		return "ok", nil
	}).DependsOn("critical_node"))

	_, err := Execute(context.Background(), dag)
	if err == nil {
		t.Fatal("expected error from critical node")
	}
	if downstreamRan.Load() {
		t.Fatal("downstream should not have run after critical failure")
	}
}

func TestNonCriticalNodeFailureDoesNotCancel(t *testing.T) {
	dag := NewDAG()

	dag.AddNode(NewNode("non_critical", func(ctx context.Context, deps map[string]any) (any, error) {
		return nil, errors.New("non-critical failure")
	}).NonCritical())

	dag.AddNode(NewNode("independent", func(ctx context.Context, deps map[string]any) (any, error) {
		return "independent_result", nil
	}))

	dag.AddNode(NewNode("depends_on_non_critical", func(ctx context.Context, deps map[string]any) (any, error) {
		// deps["non_critical"] will be nil since it failed
		val := deps["non_critical"]
		if val != nil {
			t.Errorf("expected nil dep from failed non-critical node, got %v", val)
		}
		return "downstream_ok", nil
	}).DependsOn("non_critical"))

	result, err := Execute(context.Background(), dag)
	if err != nil {
		t.Fatalf("DAG should not fail for non-critical node error: %v", err)
	}

	// independent node should have run fine
	val := result.MustGet("independent").(string)
	if val != "independent_result" {
		t.Fatalf("expected independent_result, got %s", val)
	}

	// downstream of non-critical should still run
	val = result.MustGet("depends_on_non_critical").(string)
	if val != "downstream_ok" {
		t.Fatalf("expected downstream_ok, got %s", val)
	}

	// non-critical error should be recorded
	errs := result.Errors()
	if _, ok := errs["non_critical"]; !ok {
		t.Fatal("expected non-critical error to be recorded")
	}
}

func TestMixedCriticalAndNonCritical(t *testing.T) {
	dag := NewDAG()

	// non-critical fails but doesn't cancel
	dag.AddNode(NewNode("marketing", func(ctx context.Context, deps map[string]any) (any, error) {
		return nil, errors.New("marketing service down")
	}).NonCritical())

	// critical succeeds
	dag.AddNode(NewNode("settle", func(ctx context.Context, deps map[string]any) (any, error) {
		return "settle_ok", nil
	}))

	// depends on both
	dag.AddNode(NewNode("aggregate", func(ctx context.Context, deps map[string]any) (any, error) {
		settle := deps["settle"].(string)
		marketing := deps["marketing"] // nil because it failed
		return fmt.Sprintf("settle=%s, marketing=%v", settle, marketing), nil
	}).DependsOn("settle", "marketing"))

	result, err := Execute(context.Background(), dag)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	val := result.MustGet("aggregate").(string)
	if val != "settle=settle_ok, marketing=<nil>" {
		t.Fatalf("unexpected aggregate result: %s", val)
	}

	errs := result.Errors()
	if len(errs) != 1 {
		t.Fatalf("expected 1 non-critical error, got %d", len(errs))
	}
}

func TestRetryOnFailure(t *testing.T) {
	dag := NewDAG()

	var attempts atomic.Int32
	dag.AddNode(NewNode("flaky", func(ctx context.Context, deps map[string]any) (any, error) {
		n := attempts.Add(1)
		if n < 3 {
			return nil, errors.New("temporary error")
		}
		return "success", nil
	}).WithRetry(RetryPolicy{
		MaxRetries:    3,
		InitialDelay:  10 * time.Millisecond,
		BackoffFactor: 1.0,
	}))

	result, err := Execute(context.Background(), dag)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	val := result.MustGet("flaky").(string)
	if val != "success" {
		t.Fatalf("expected success, got %s", val)
	}
	if attempts.Load() != 3 {
		t.Fatalf("expected 3 attempts, got %d", attempts.Load())
	}
}

func TestRetryExhausted(t *testing.T) {
	dag := NewDAG()

	dag.AddNode(NewNode("always_fail", func(ctx context.Context, deps map[string]any) (any, error) {
		return nil, errors.New("permanent error")
	}).WithRetry(RetryPolicy{
		MaxRetries:   2,
		InitialDelay: 10 * time.Millisecond,
	}))

	_, err := Execute(context.Background(), dag)
	if err == nil {
		t.Fatal("expected error after retries exhausted")
	}
}

func TestEndpointRotation(t *testing.T) {
	dag := NewDAG()

	var hitEndpoints []string
	var mu sync.Mutex

	dag.AddNode(NewEndpointNode("service", func(ctx context.Context, deps map[string]any, endpoint string) (any, error) {
		mu.Lock()
		hitEndpoints = append(hitEndpoints, endpoint)
		mu.Unlock()

		// First two endpoints "fail" (simulating timeout)
		if endpoint == "host1:8080" || endpoint == "host2:8080" {
			return nil, errors.New("connection timeout")
		}
		return "ok from " + endpoint, nil
	}, RetryPolicy{
		MaxRetries:    3,
		InitialDelay:  10 * time.Millisecond,
		BackoffFactor: 1.0,
		Endpoints:     []string{"host1:8080", "host2:8080", "host3:8080"},
	}))

	result, err := Execute(context.Background(), dag)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	val := result.MustGet("service").(string)
	if val != "ok from host3:8080" {
		t.Fatalf("expected ok from host3:8080, got %s", val)
	}

	mu.Lock()
	defer mu.Unlock()
	if len(hitEndpoints) != 3 {
		t.Fatalf("expected 3 endpoint hits, got %d: %v", len(hitEndpoints), hitEndpoints)
	}
	// Should have rotated: host1 -> host2 -> host3
	if hitEndpoints[0] != "host1:8080" || hitEndpoints[1] != "host2:8080" || hitEndpoints[2] != "host3:8080" {
		t.Fatalf("unexpected endpoint rotation: %v", hitEndpoints)
	}
}

func TestRetryRespectsIsRetryable(t *testing.T) {
	dag := NewDAG()

	var attempts atomic.Int32
	dag.AddNode(NewNode("non_retryable", func(ctx context.Context, deps map[string]any) (any, error) {
		attempts.Add(1)
		return nil, errors.New("4xx: bad request")
	}).WithRetry(RetryPolicy{
		MaxRetries:   3,
		InitialDelay: 10 * time.Millisecond,
		IsRetryable: func(err error) bool {
			// Only retry 5xx errors
			return !errors.Is(err, err) || err.Error() != "4xx: bad request"
		},
	}))

	_, err := Execute(context.Background(), dag)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestRetryRespectsContextCancellation(t *testing.T) {
	dag := NewDAG()

	dag.AddNode(NewNode("slow_retry", func(ctx context.Context, deps map[string]any) (any, error) {
		return nil, errors.New("keep failing")
	}).WithRetry(RetryPolicy{
		MaxRetries:   10,
		InitialDelay: 500 * time.Millisecond,
	}))

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	_, err := Execute(ctx, dag)
	if err == nil {
		t.Fatal("expected timeout error")
	}
}
