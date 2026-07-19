package predicate_test

import (
	"fmt"
	"testing"

	"github.com/amirali-amirifar/conform/predicate"
	"github.com/stretchr/testify/require"
)

// Node[T] must be satisfiable by a type parameter constrained to conform's
// integer types (a subset of cmp.Ordered) as well as a defined ~int type.
type myInt int

var (
	_ predicate.Node[int]   = predicate.Cmp[int]{}
	_ predicate.Node[myInt] = predicate.NewCmp(predicate.Ge, myInt(1))
)

// even is a user-defined constraint living outside the predicate package. It
// implements predicate.Node[int] with no help from the library, proving a third
// party can add a validator and compose it with the built-in nodes.
type even struct{}

func (even) Validate(v int) []predicate.Diagnostic {
	if v%2 == 0 {
		return nil
	}
	return []predicate.Diagnostic{{
		Predicate: "even",
		Message:   fmt.Sprintf("got %d but need an even number", v),
	}}
}

func (even) String() string { return "x is even" }

func TestCustomNode(t *testing.T) {
	r := require.New(t)
	// A built-in bound AND a caller-supplied rule, in the same tree.
	tree := predicate.NewAnd(predicate.NewCmp(predicate.Ge, 0), even{})

	r.Empty(predicate.Validate(tree, 4))   // satisfies both
	r.Len(predicate.Validate(tree, 3), 1)  // even fails
	r.Len(predicate.Validate(tree, -3), 2) // both fail
	r.Equal("(x >= 0 AND x is even)", predicate.String(tree))
}

// Or short-circuits on the first satisfied branch, and the empty-tree
// identities differ: empty Or is unsatisfiable, empty And is vacuously true.
// (Per-leaf diagnostics for And/Cmp/In/Not are covered in validate_test.go.)
func TestValidateOr(t *testing.T) {
	r := require.New(t)
	// v < 0 OR v > 100
	tree := predicate.NewOr(
		predicate.NewCmp(predicate.Lt, 0),
		predicate.NewCmp(predicate.Gt, 100),
	)

	r.Empty(predicate.Validate(tree, 200)) // second branch matches
	r.Empty(predicate.Validate(tree, -5))  // first branch matches

	// A total Or failure is one diagnostic that groups both branches, not two
	// flat entries that would read as separate requirements.
	diags := predicate.Validate(tree, 50)
	r.Len(diags, 1)
	r.Equal("Or", diags[0].Predicate)
	r.Len(diags[0].Children, 2)

	// Empty disjunction is unsatisfiable (identity of OR), unlike empty And.
	r.NotEmpty(predicate.Validate(predicate.NewOr[int](), 1))
	r.Empty(predicate.Validate(predicate.NewAnd[int](), 1))
}

// String renders a tree that nests every node type, exercising the recursive
// composition of the per-node String methods.
func TestString(t *testing.T) {
	r := require.New(t)
	// 1 <= x <= 100 AND (x in {10,20,30} OR x != 50)
	tree := predicate.NewAnd(
		predicate.NewCmp(predicate.Ge, 1),
		predicate.NewCmp(predicate.Le, 100),
		predicate.NewOr(
			predicate.NewIn(10, 20, 30),
			predicate.NewNot(predicate.NewCmp(predicate.Eq, 50)),
		),
	)
	r.Equal(
		"(x >= 1 AND x <= 100 AND (x in [10 20 30] OR NOT x == 50))",
		predicate.String(tree),
	)
}
