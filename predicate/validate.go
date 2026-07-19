package predicate

import "cmp"

// Validate checks v against the constraint tree n and returns one [Diagnostic]
// per failed leaf. An empty (nil) result means v satisfies the tree.
//
// It is a nil-safe wrapper over [Node.Validate] so a caller can hold an
// optional (possibly nil) root; a nil node passes vacuously. Dispatch is by
// method, so a tree may freely mix built-in and user-defined nodes.
func Validate[T cmp.Ordered](n Node[T], v T) []Diagnostic {
	if n == nil {
		return nil
	}
	return n.Validate(v)
}
