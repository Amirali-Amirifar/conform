package conform

import "github.com/amirali-amirifar/conform/predicate"

type spec[T IntType] struct {
	root  predicate.Node[T]
	built bool
}

func NewSpec[T IntType](root predicate.Node[T]) spec[T] {
	return spec[T]{
		root:  root,
		built: true,
	}
}

// validate returns one diagnostic per failed constraint, or none when v passes.
// An unbuilt spec (a zero-value box) rejects every value.
func (s spec[T]) validate(v T) []predicate.Diagnostic {
	if !s.built {
		return []predicate.Diagnostic{{
			Predicate: "Spec",
			Message:   "conform: Spec is not initialized; construct it with NewInt",
		}}
	}
	return predicate.Validate(s.root, v)
}
