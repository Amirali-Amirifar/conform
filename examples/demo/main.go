// Command demo exercises conform end to end: a config struct validated as it is
// decoded, and a large hand-built predicate tree validated directly so every
// node kind (And, Or, Not, Cmp, In, and a custom node) shows up in one run.
//
//	go run ./examples/demo
package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/amirali-amirifar/conform"
	"github.com/amirali-amirifar/conform/predicate"
)

// Even is a user-defined constraint, living entirely outside the library. It
// implements predicate.Node[int] and composes with the built-in nodes.
type Even struct{}

func (Even) Validate(v int) []predicate.Diagnostic {
	if v%2 == 0 {
		return nil
	}
	return []predicate.Diagnostic{{
		Predicate: "Even",
		Message:   fmt.Sprintf("got %d, need an even number", v),
	}}
}

func (Even) String() string { return "x is even" }

// ServerConfig is validated field-by-field as it is decoded.
type ServerConfig struct {
	Workers  conform.Int[int] `json:"workers"`
	LogEvery conform.Int[int] `json:"log_every"`
}

func NewServerConfig() ServerConfig {
	return ServerConfig{
		Workers:  conform.NewInt(predicate.NewAnd(conform.Min(1), conform.Max(64))),
		LogEvery: conform.NewInt(conform.In(1, 10, 100, 1000)),
	}
}

func main() {
	demoConfig()
	demoTree()
}

func demoConfig() {
	section("conform config — validation happens during json.Unmarshal")

	for _, input := range []string{
		`{"workers": 8, "log_every": 100}`, // valid
		`{"workers": 200, "log_every": 7}`, // both out of range
	} {
		cfg := NewServerConfig()
		fmt.Println("input:", input)
		if err := json.Unmarshal([]byte(input), &cfg); err != nil {
			fmt.Printf("  rejected:\n%s\n\n", indentLines(err.Error(), "    "))
			continue
		}
		fmt.Printf("  accepted: workers=%d log_every=%d\n\n", cfg.Workers.Value(), cfg.LogEvery.Value())
	}
}

func demoTree() {
	section("predicate tree — one big tree covering every node kind")

	// 1 <= x <= 100 AND x in {10,20,30} AND NOT x == 251
	//   AND (x < 0 OR x > 1000)      <- disjunction, both branches fail
	//   AND (x >= 5 OR x in {1,2,3}) <- disjunction, first branch passes
	//   AND (x is even)              <- custom node
	tree := predicate.NewAnd(
		predicate.NewCmp(predicate.Ge, 1),
		predicate.NewCmp(predicate.Le, 100),
		predicate.NewIn(10, 20, 30),
		predicate.NewNot(predicate.NewCmp(predicate.Eq, 251)),
		predicate.NewOr(
			predicate.NewCmp(predicate.Lt, 0),
			predicate.NewCmp(predicate.Gt, 1000),
		),
		predicate.NewOr(
			predicate.NewCmp(predicate.Ge, 5),
			predicate.NewIn(1, 2, 3),
		),
		Even{},
	)

	fmt.Println("rendered:", predicate.String(tree))
	fmt.Println("\nstructure (walked by type-switch over the exported node types):")
	printTree(tree, "  ")

	const v = 251
	fmt.Printf("\nvalidating value = %d\n", v)
	diags := predicate.Validate(tree, v)
	if len(diags) == 0 {
		fmt.Println("  valid!")
		return
	}
	fmt.Printf("\ndiagnostics tree (%d top-level, Or groups its branches):\n", len(diags))
	printDiagnostics(diags, "  ")
}

// printTree walks the predicate tree by type-switching over the exported node
// types — the same pattern an outside code generator would use. An unrecognized
// node (here, the custom Even) falls through to its String form.
func printTree(n predicate.Node[int], indent string) {
	switch t := n.(type) {
	case predicate.And[int]:
		fmt.Printf("%sAND\n", indent)
		for _, c := range t.Preds {
			printTree(c, indent+"  ")
		}
	case predicate.Or[int]:
		fmt.Printf("%sOR\n", indent)
		for _, c := range t.Preds {
			printTree(c, indent+"  ")
		}
	case predicate.Not[int]:
		fmt.Printf("%sNOT\n", indent)
		printTree(t.Pred, indent+"  ")
	case predicate.Cmp[int]:
		fmt.Printf("%sCmp: x %s %d\n", indent, t.Op, t.Val)
	case predicate.In[int]:
		fmt.Printf("%sIn: x in %v\n", indent, t.Allowed)
	default:
		fmt.Printf("%s%s (custom node)\n", indent, n)
	}
}

func printDiagnostics(diags []predicate.Diagnostic, indent string) {
	for _, d := range diags {
		fmt.Printf("%s- [%s] %s\n", indent, d.Predicate, d.Message)
		if len(d.Children) > 0 {
			printDiagnostics(d.Children, indent+"    ")
		}
	}
}

func section(title string) {
	fmt.Printf("\n=== %s ===\n\n", title)
}

func indentLines(s, indent string) string {
	return indent + strings.ReplaceAll(s, "\n", "\n"+indent)
}
