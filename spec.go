package conform

import "fmt"

type spec[T IntType] struct {
	min        T
	max        T
	configured bool
}

func NewSpec[T IntType](vMin, vMax T) spec[T] {
	if vMin > vMax {
		panic("conform: vMin must not be greater than vMax")
	}
	return spec[T]{min: vMin, max: vMax, configured: true}
}

func (s spec[T]) validate(v T) error {
	if !s.configured {
		return fmt.Errorf("conform: Int is not initialized; construct it with NewInt")
	}
	if v < s.min {
		return fmt.Errorf("got %d but need at least %d", v, s.min)
	}
	if v > s.max {
		return fmt.Errorf("got %d but need at most %d", v, s.max)
	}
	return nil
}
