// Package predicate is a small, dependency-free kernel for boolean constraint
// trees over a single ordered value.
//
// A [Node] is reified, walkable data: unlike an opaque func(v T) error, a tree
// of nodes can be inspected, rendered, and composed. Validation is a method on
// the node, so anyone can add a constraint by implementing [Node] — a custom
// rule composes with the built-in [And], [Or], [Not], [Cmp] and [In] just like
// any other node.
//
// A future interpreter that needs the tree's structure (say, a SQL code
// generator) type-switches over the exported node types and decides for itself
// how to treat a node it does not recognize. The tree describes constraints on
// a single value of type T; a heterogeneous predicate across many
// differently-typed fields is a separate structure that reuses these patterns
// rather than this exact type.
package predicate

import (
	"cmp"
	"fmt"
)

// Node is one node of a constraint tree over a value of type T.
//
// Implementations are plain structs with exported fields (see [And], [Or],
// [Not], [Cmp], [In]); the tree is inert data apart from the two behaviours
// every node supports: validating a value and rendering itself. Implement Node
// to contribute a custom constraint.
type Node[T cmp.Ordered] interface {
	fmt.Stringer

	// Validate reports one Diagnostic per way v fails this node's constraint;
	// an empty result means v satisfies it.
	Validate(v T) []Diagnostic
}

// Diagnostic explains one way a value failed a constraint.
type Diagnostic struct {
	// Predicate is the name of the predicate node that produced this
	// diagnostic (for example "Cmp" or "In"). It lets a consumer react to a
	// failure without parsing Message.
	Predicate string
	// Message is the human-readable reason the value failed.
	Message string
	// Children holds the sub-diagnostics of a composite predicate, mirroring
	// the walk down the tree: an Or groups the failures of its alternatives
	// here. It is nil for leaves and for conjunctions (an And flattens its
	// children into the returned slice, since every one of them is required).
	Children []Diagnostic
}

// String renders the diagnostic as "Predicate: Message", or just Message when
// Predicate is empty.
func (d Diagnostic) String() string {
	if d.Predicate == "" {
		return d.Message
	}
	return d.Predicate + ": " + d.Message
}
