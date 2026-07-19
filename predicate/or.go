package predicate

import (
	"cmp"
	"fmt"
)

// Or is satisfied when at least one child is satisfied.
type Or[T cmp.Ordered] struct {
	Preds []Node[T]
}

// NewOr builds an [Or] over preds.
func NewOr[T cmp.Ordered](preds ...Node[T]) Or[T] {
	return Or[T]{Preds: preds}
}

// Validate passes as soon as one child passes. An empty Or is the identity of
// OR: nothing satisfies it.
func (n Or[T]) Validate(v T) []Diagnostic {
	if len(n.Preds) == 0 {
		return []Diagnostic{{
			Predicate: "Or",
			Message:   fmt.Sprintf("got %v but no alternative is allowed", v),
		}}
	}
	var children []Diagnostic
	for _, p := range n.Preds {
		d := Validate(p, v)
		if len(d) == 0 {
			return nil // one branch satisfied is enough
		}
		children = append(children, d...)
	}
	// No branch matched: report one Or failure that groups the alternatives'
	// diagnostics, so the disjunction isn't mistaken for a list of requirements.
	return []Diagnostic{{
		Predicate: "Or",
		Message:   fmt.Sprintf("got %v but matched none of %s", v, n.String()),
		Children:  children,
	}}
}

func (n Or[T]) String() string { return join(n.Preds, " OR ") }
