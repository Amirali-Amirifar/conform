package conform_test

import (
	"encoding/json"
	"testing"

	"github.com/amirali-amirifar/conform"
	"gopkg.in/yaml.v3"
)

func TestUnmarshal(t *testing.T) {
	type S struct {
		Val conform.Int[int8] `json:"val"`
	}

	t.Run("JSON valid", func(t *testing.T) {
		cfg := S{Val: conform.NewInt[int8](1, 100)}
		if err := json.Unmarshal([]byte(`{"val": 89}`), &cfg); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if v := cfg.Val.Value(); v != 89 {
			t.Fatalf("got %v, want 89", v)
		}
	})

	t.Run("JSON out of range", func(t *testing.T) {
		cfg := S{Val: conform.NewInt[int8](1, 100)}
		if err := json.Unmarshal([]byte(`{"val": 120}`), &cfg); err == nil {
			t.Fatal("expected error for out-of-range value, got nil")
		}
	})

	t.Run("YAML valid", func(t *testing.T) {
		cfg := S{Val: conform.NewInt[int8](1, 100)}

		if err := yaml.Unmarshal([]byte("val: 89"), &cfg); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if v := cfg.Val.Value(); v != 89 {
			t.Fatalf("got %v, want 89", v)
		}
	})

	t.Run("YAML out of range", func(t *testing.T) {
		cfg := struct {
			Val conform.Int[int8] `yaml:"val"`
		}{Val: conform.NewInt[int8](1, 100)}

		if err := yaml.Unmarshal([]byte("val: 120"), &cfg); err == nil {
			t.Fatal("expected error for out-of-range value, got nil")
		}
	})

	t.Run("zero value rejects unmarshal", func(t *testing.T) {
		var cfg struct {
			Val conform.Int[int8] `json:"val"`
		}

		if err := json.Unmarshal([]byte(`{"val": 42}`), &cfg); err == nil {
			t.Fatal("expected error for uninitialized Int, got nil")
		}
		if cfg.Val.IsValid() {
			t.Fatal("uninitialized Int must not report a valid value")
		}
	})
}
