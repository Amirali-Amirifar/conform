package predicate_test

import (
	"testing"

	"github.com/amirali-amirifar/conform/predicate"
	"github.com/stretchr/testify/require"
)

// TestValidateDiagnostics validates one value against a nested tree that it
// fails in several independent ways, showing that Validate collects a
// diagnostic per failing leaf (not just the first) and that the wording is the
// contract the README quotes. The passing leaf (Ge 1) contributes nothing.
func TestValidateDiagnostics(t *testing.T) {
	r := require.New(t)

	tree := predicate.NewAnd(
		predicate.NewCmp(predicate.Ge, 1),                     // passes at 200
		predicate.NewCmp(predicate.Le, 100),                   // fails: too big
		predicate.NewIn(10, 20, 30),                           // fails: not allow-listed
		predicate.NewNot(predicate.NewCmp(predicate.Eq, 200)), // fails: forbidden value
	)

	r.Equal([]predicate.Diagnostic{
		{Predicate: "Cmp", Message: "got 200 but need at most 100"},
		{Predicate: "In", Message: "got 200 but need one of [10 20 30]"},
		{Predicate: "Not", Message: "got 200 but expected NOT (x == 200)"},
	}, predicate.Validate(tree, 200))

	// A value that satisfies every leaf yields no diagnostics.
	r.Nil(predicate.Validate(tree, 20))
}
