package predicate

import (
	"cmp"
	"strings"
)

// String renders n as an infix expression, e.g. "(x >= 1 AND x <= 64)". It is
// a nil-safe wrapper over [Node.String]; a nil node renders as "true".
func String[T cmp.Ordered](n Node[T]) string {
	if n == nil {
		return "true"
	}
	return n.String()
}

// join renders a parenthesized list of children separated by sep, used by the
// [And] and [Or] String methods.
func join[T cmp.Ordered](preds []Node[T], sep string) string {
	parts := make([]string, len(preds))
	for i, p := range preds {
		parts[i] = String(p)
	}
	return "(" + strings.Join(parts, sep) + ")"
}
