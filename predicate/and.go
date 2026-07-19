package predicate

import "cmp"

// And is satisfied when every child is satisfied.
type And[T cmp.Ordered] struct {
	Preds []Node[T]
}

// NewAnd builds an [And] over preds.
func NewAnd[T cmp.Ordered](preds ...Node[T]) And[T] {
	return And[T]{Preds: preds}
}

// Validate reports every child's diagnostics; an empty And is vacuously
// satisfied.
func (n And[T]) Validate(v T) []Diagnostic {
	var diags []Diagnostic
	for _, p := range n.Preds {
		diags = append(diags, Validate(p, v)...)
	}
	return diags
}

func (n And[T]) String() string { return join(n.Preds, " AND ") }
