package conform

import "errors"

type spec[T IntType] struct {
	rules []Rule[T]
}

func NewSpec[T IntType](rules ...Rule[T]) spec[T] {
	if rules == nil {
		rules = []Rule[T]{}
	}
	return spec[T]{rules: rules}
}

func (s spec[T]) validate(v T) error {
	if s.rules == nil {
		return errors.New("conform: Spec is not initialized; construct it with NewInt")
	}
	var errs []error
	for _, r := range s.rules {
		if err := r(v); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...) // nil when errs is empty
}
