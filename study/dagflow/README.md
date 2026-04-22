# dagflow

A lightweight DAG-based async orchestration framework for Go.

## Features

- DAG dependency graph with cycle detection
- Topological execution order with maximum parallelism
- Dependency result passing between nodes
- Goroutine pool for concurrency control
- Context timeout and cancellation propagation
- Node failure cancels the entire DAG via errgroup

## Install

```bash
go get github.com/zhoujie0420/dagflow
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "github.com/zhoujie0420/dagflow"
)

func main() {
    dag := dagflow.NewDAG()

    dag.AddNode(dagflow.NewNode("settle", func(ctx context.Context, deps map[string]any) (any, error) {
        return "settle_result", nil
    }))

    dag.AddNode(dagflow.NewNode("risk", func(ctx context.Context, deps map[string]any) (any, error) {
        return "risk_passed", nil
    }))

    dag.AddNode(dagflow.NewNode("card", func(ctx context.Context, deps map[string]any) (any, error) {
        return map[string]string{"org": "VISA"}, nil
    }))

    // route depends on card result
    dag.AddNode(dagflow.NewNode("route", func(ctx context.Context, deps map[string]any) (any, error) {
        card := deps["card"].(map[string]string)
        return fmt.Sprintf("route for %s", card["org"]), nil
    }).DependsOn("card"))

    // aggregate depends on all
    dag.AddNode(dagflow.NewNode("aggregate", func(ctx context.Context, deps map[string]any) (any, error) {
        return fmt.Sprintf("settle=%v, risk=%v, route=%v",
            deps["settle"], deps["risk"], deps["route"]), nil
    }).DependsOn("settle", "risk", "route"))

    result, err := dagflow.Execute(context.Background(), dag, dagflow.WithWorkers(3))
    if err != nil {
        panic(err)
    }

    fmt.Println(result.MustGet("aggregate"))
}
```

## API

### Node

```go
// Create a node
node := dagflow.NewNode("name", func(ctx context.Context, deps map[string]any) (any, error) {
    // deps contains results from dependency nodes
    return "result", nil
})

// Declare dependencies
node.DependsOn("dep1", "dep2")
```

### DAG

```go
dag := dagflow.NewDAG()
dag.AddNode(node)       // Add node, returns error if duplicate name
dag.Validate()          // Check for cycles and missing dependencies
```

### Execute

```go
// Execute with default concurrency (one goroutine per node)
result, err := dagflow.Execute(ctx, dag)

// Execute with limited concurrency via goroutine pool
result, err := dagflow.Execute(ctx, dag, dagflow.WithWorkers(5))

// Read results
val, ok := result.Get("node_name")
val := result.MustGet("node_name") // panics if not found
```

## License

MIT
