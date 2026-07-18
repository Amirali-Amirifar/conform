package conform

import (
	"encoding/json"
	"fmt"

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

func NewInt[T IntType](rules ...Rule[T]) Int[T] {
	return Int[T]{
		spec: NewSpec(rules...),
	}
}

func (i Int[T]) Parse(v T) (Int[T], error) {
	if err := i.spec.validate(v); err != nil {
		return Int[T]{}, err
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
	parsed, err := i.Parse(raw)
	if err != nil {
		return err
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
	parsed, err := i.Parse(raw)
	if err != nil {
		return err
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
