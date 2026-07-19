package predicate

import (
	"cmp"
	"fmt"
)

// Op is a comparison operator used by a [Cmp] leaf.
type Op int

// Operators start at 1 so the zero value is not a valid Op: an uninitialized
// Cmp is meaningless and is caught loudly rather than behaving like ==.
const (
	Eq Op = iota + 1 // ==
	Ne               // !=
	Lt               // <
	Le               // <=
	Gt               // >
	Ge               // >=
)

// String returns the operator's symbol, e.g. ">=".
func (o Op) String() string {
	switch o {
	case Eq:
		return "=="
	case Ne:
		return "!="
	case Lt:
		return "<"
	case Le:
		return "<="
	case Gt:
		return ">"
	case Ge:
		return ">="
	default:
		panic(fmt.Sprintf("predicate: invalid Op(%d)", int(o)))
	}
}

// phrase returns a human-readable form of the operator for diagnostics,
// e.g. "at least" for Ge. String returns the symbol; phrase reads in prose.
func (o Op) phrase() string {
	switch o {
	case Eq:
		return "exactly"
	case Ne:
		return "not"
	case Lt:
		return "less than"
	case Le:
		return "at most"
	case Gt:
		return "greater than"
	case Ge:
		return "at least"
	default:
		panic(fmt.Sprintf("predicate: invalid Op(%d)", int(o)))
	}
}

// evalOp reports whether "a op b" holds.
func evalOp[T cmp.Ordered](o Op, a, b T) bool {
	switch o {
	case Eq:
		return a == b
	case Ne:
		return a != b
	case Lt:
		return a < b
	case Le:
		return a <= b
	case Gt:
		return a > b
	case Ge:
		return a >= b
	default:
		panic(fmt.Sprintf("predicate: invalid Op(%d)", int(o)))
	}
}
