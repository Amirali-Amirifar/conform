package predicate

// Predicates

type AndPredicate[T any] struct {
	preds []Predicate[T]
}

func And[T any](preds ...Predicate[T]) AndPredicate[T] {
	return AndPredicate[T]{
		preds: preds,
	}
}

func (a AndPredicate[T]) Name() string {
	return "AND"
}

func (a AndPredicate[T]) Evaluate(in T) bool {
	for _, p := range a.preds {
		if diags := p.Evaluate(in); diags != nil {
			return false
		}
	}
	return true
}
