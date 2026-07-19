package predicate

import (
	"cmp"
	"fmt"
	"slices"
)

// In is a leaf satisfied when the tested value equals one of Allowed.
type In[T cmp.Ordered] struct {
	Allowed []T
}

// NewIn builds an [In] leaf from an allow-list.
func NewIn[T cmp.Ordered](allowed ...T) In[T] {
	return In[T]{Allowed: allowed}
}

// Validate checks membership of v in Allowed.
func (n In[T]) Validate(v T) []Diagnostic {
	if slices.Contains(n.Allowed, v) {
		return nil
	}
	return []Diagnostic{{
		Predicate: "In",
		Message:   fmt.Sprintf("got %v but need one of %v", v, n.Allowed),
	}}
}

func (n In[T]) String() string { return fmt.Sprintf("x in %v", n.Allowed) }
