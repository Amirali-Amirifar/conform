package conform

import "fmt"

// Rule reports whether v satisfies a constraint. A nil error means it passes.
type Rule[T any] func(v T) error

func Min[T IntType](min T) Rule[T] {
	return func(v T) error {
		if v < min {
			return fmt.Errorf("got %d but need at least %d", v, min)
		}
		return nil
	}
}

func Max[T IntType](max T) Rule[T] {
	return func(v T) error {
		if v > max {
			return fmt.Errorf("got %d but need at most %d", v, max)
		}
		return nil
	}
}

// In restricts v to an allow-list of values.
func In[T IntType](allowed ...T) Rule[T] {
	return func(v T) error {
		for _, a := range allowed {
			if v == a {
				return nil
			}
		}
		return fmt.Errorf("got %d but need one of %v", v, allowed)
	}
}
