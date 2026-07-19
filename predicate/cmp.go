package predicate

import (
	"cmp"
	"fmt"
)

// Cmp is a leaf satisfied when the tested value compares to Val under Op,
// i.e. "value Op Val" holds.
type Cmp[T cmp.Ordered] struct {
	Op  Op
	Val T
}

// NewCmp builds a [Cmp] leaf.
func NewCmp[T cmp.Ordered](op Op, val T) Cmp[T] {
	return Cmp[T]{Op: op, Val: val}
}

// Validate compares v against Val under Op.
func (n Cmp[T]) Validate(v T) []Diagnostic {
	if !evalOp(n.Op, v, n.Val) {
		return []Diagnostic{{
			Predicate: "Cmp",
			Message:   fmt.Sprintf("got %v but need %s %v", v, n.Op.phrase(), n.Val),
		}}
	}
	return nil
}

func (n Cmp[T]) String() string { return fmt.Sprintf("x %s %v", n.Op, n.Val) }
