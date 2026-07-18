package predicate

import "fmt"

type Diagnostic string

type Predicate[T any] interface {
	fmt.Stringer

	Evaluate(T) []Diagnostic
}
