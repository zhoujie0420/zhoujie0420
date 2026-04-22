package dagflow

import (
	"fmt"
	"sync"
)

// DAG represents a directed acyclic graph of nodes.
type DAG struct {
	nodes map[string]*Node
	order []string // insertion order for deterministic iteration
	mu    sync.RWMutex
}

// NewDAG creates an empty DAG.
func NewDAG() *DAG {
	return &DAG{
		nodes: make(map[string]*Node),
	}
}

// AddNode adds a node to the DAG. Returns error if a node with the same name already exists.
func (d *DAG) AddNode(node *Node) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if _, exists := d.nodes[node.name]; exists {
		return fmt.Errorf("dagflow: duplicate node name %q", node.name)
	}
	d.nodes[node.name] = node
	d.order = append(d.order, node.name)
	return nil
}

// Validate checks that:
// 1. All dependency references point to existing nodes.
// 2. The graph has no cycles.
func (d *DAG) Validate() error {
	d.mu.RLock()
	defer d.mu.RUnlock()

	// Check all deps exist
	for _, node := range d.nodes {
		for _, dep := range node.deps {
			if _, ok := d.nodes[dep]; !ok {
				return fmt.Errorf("dagflow: node %q depends on unknown node %q", node.name, dep)
			}
		}
	}

	// Cycle detection using Kahn's algorithm
	inDegree := make(map[string]int, len(d.nodes))
	downstream := make(map[string][]string, len(d.nodes))

	for name := range d.nodes {
		inDegree[name] = 0
	}
	for _, node := range d.nodes {
		for _, dep := range node.deps {
			downstream[dep] = append(downstream[dep], node.name)
			inDegree[node.name]++
		}
	}

	queue := make([]string, 0)
	for name, deg := range inDegree {
		if deg == 0 {
			queue = append(queue, name)
		}
	}

	visited := 0
	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]
		visited++
		for _, next := range downstream[curr] {
			inDegree[next]--
			if inDegree[next] == 0 {
				queue = append(queue, next)
			}
		}
	}

	if visited != len(d.nodes) {
		return fmt.Errorf("dagflow: cycle detected in graph")
	}

	return nil
}
