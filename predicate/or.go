package predicate

// Predicates

type OrPredicate[T any] struct {
	preds []Predicate[T]
}

func Or[T any](preds ...Predicate[T]) OrPredicate[T] {
	return OrPredicate[T]{
		preds: preds,
	}
}

func (a OrPredicate[T]) Name() string {
	return "AND"
}

func (a OrPredicate[T]) Evaluate(in T) bool {
	for _, p := range a.preds {
		if diags := p.Evaluate(in); diags == nil {
			return true
		}
	}
	return false
}
