package dagflow

import "context"

// NodeFunc is the function signature for a node's execution logic.
// deps contains the results from all dependency nodes, keyed by node name.
type NodeFunc func(ctx context.Context, deps map[string]interface{}) (interface{}, error)

// Node represents a single unit of work in the DAG.
type Node struct {
	name     string
	deps     []string
	fn       NodeFunc
	critical bool         // if true, failure cancels the entire DAG
	retry    *RetryPolicy // optional retry policy
}

// NewNode creates a new node with the given name and execution function.
// By default, nodes are critical — failure cancels the entire DAG.
func NewNode(name string, fn NodeFunc) *Node {
	return &Node{
		name:     name,
		fn:       fn,
		critical: true,
	}
}

// DependsOn declares that this node depends on the given node names.
// The results of those nodes will be available in the deps map when this node executes.
func (n *Node) DependsOn(names ...string) *Node {
	n.deps = append(n.deps, names...)
	return n
}

// NonCritical marks this node as non-critical.
// If a non-critical node fails, its error is logged but does not cancel the DAG.
// Downstream nodes will receive nil for this node's result in their deps map.
func (n *Node) NonCritical() *Node {
	n.critical = false
	return n
}

// WithRetry configures retry policy for this node.
// On failure, the node function will be retried up to policy.MaxRetries times.
func (n *Node) WithRetry(policy RetryPolicy) *Node {
	n.retry = &policy
	return n
}

// NewEndpointNode creates a node that rotates through endpoints on retry.
// The EndpointFunc receives the current endpoint address as a parameter.
// Requires a RetryPolicy with Endpoints configured.
func NewEndpointNode(name string, fn EndpointFunc, policy RetryPolicy) *Node {
	return &Node{
		name:     name,
		fn:       withEndpointRetry(&policy, fn),
		critical: true,
	}
}

// Name returns the node's name.
func (n *Node) Name() string {
	return n.name
}
