package conform

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/amirali-amirifar/conform/predicate"
	"gopkg.in/yaml.v3"
)

type IntType interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

type Int[T IntType] struct {
	spec  spec[T]
	value T
	valid bool
}

func NewInt[T IntType](root predicate.Node[T]) Int[T] {
	return Int[T]{
		spec: NewSpec(root),
	}
}

// Parse checks v against the spec and returns the populated Int, or the
// diagnostics explaining why v was rejected. A nil/empty result means success.
func (i Int[T]) Parse(v T) (Int[T], []predicate.Diagnostic) {
	if diags := i.spec.validate(v); len(diags) > 0 {
		return Int[T]{}, diags
	}
	i.value = v
	i.valid = true

	return i, nil
}

func (i Int[T]) Value() T {
	if !i.valid {
		panic("conform: Int has no valid value; construct with NewInt and Parse/unmarshal first")
	}
	return i.value
}

func (i Int[T]) IsValid() bool {
	return i.valid
}

func (i *Int[T]) UnmarshalJSON(data []byte) error {
	var raw T
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	parsed, diags := i.Parse(raw)
	if len(diags) > 0 {
		errs := make([]error, len(diags))
		for j, d := range diags {
			errs[j] = errors.New(d.Message)
		}
		return errors.Join(errs...)
	}
	*i = parsed
	return nil
}

func (i Int[T]) MarshalJSON() ([]byte, error) {
	if !i.valid {
		return nil, fmt.Errorf("conform: cannot marshal Int with no valid value")
	}
	return json.Marshal(i.value)
}

func (i *Int[T]) UnmarshalYAML(value *yaml.Node) error {
	var raw T
	if err := value.Decode(&raw); err != nil {
		return err
	}
	parsed, diags := i.Parse(raw)
	if len(diags) > 0 {
		errs := make([]error, len(diags))
		for j, d := range diags {
			errs[j] = errors.New(d.Message)
		}
		return errors.Join(errs...)
	}
	*i = parsed
	return nil
}

func (i Int[T]) MarshalYAML() (any, error) {
	if !i.valid {
		return nil, fmt.Errorf("conform: cannot marshal Int with no valid value")
	}
	return i.value, nil
}
