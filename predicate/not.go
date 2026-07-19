package predicate

import (
	"cmp"
	"fmt"
)

// Not is satisfied when its child is not satisfied.
type Not[T cmp.Ordered] struct {
	Pred Node[T]
}

// NewNot negates pred.
func NewNot[T cmp.Ordered](pred Node[T]) Not[T] {
	return Not[T]{Pred: pred}
}

// Validate fails when the child passes, and passes when the child fails.
func (n Not[T]) Validate(v T) []Diagnostic {
	if len(Validate(n.Pred, v)) == 0 {
		return []Diagnostic{{
			Predicate: "Not",
			Message:   fmt.Sprintf("got %v but expected NOT (%s)", v, String(n.Pred)),
		}}
	}
	return nil
}

func (n Not[T]) String() string { return "NOT " + String(n.Pred) }
