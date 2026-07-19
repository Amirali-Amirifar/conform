package conform_test

import (
	"encoding/json"
	"fmt"

	"github.com/amirali-amirifar/conform"
	"github.com/amirali-amirifar/conform/predicate"
)

// EngineConfig is validated as it is decoded: every field carries its own rules,
// so a successful json.Unmarshal is proof the whole config is in range.
type EngineConfig struct {
	Workers  conform.Int[int] `json:"workers"`
	LogEvery conform.Int[int] `json:"log_every"`
}

// NewEngineConfig returns a config with its validation rules pre-populated.
// Decoding into the result enforces those rules
func NewEngineConfig() EngineConfig {
	return EngineConfig{
		Workers:  conform.NewInt(predicate.NewAnd(conform.Min(1), conform.Max(64))),
		LogEvery: conform.NewInt(conform.In(1, 10, 100, 1000)),
	}
}

// Example decodes a valid config and reads the checked values.
func Example() {
	cfg := NewEngineConfig()

	const input = `{"port": 8080, "workers": 8, "log_every": 100}`
	if err := json.Unmarshal([]byte(input), &cfg); err != nil {
		fmt.Println("invalid config:", err)
		return
	}

	fmt.Println("workers:", cfg.Workers.Value())
	fmt.Println("log every:", cfg.LogEvery.Value())
	// Output:
	// workers: 8
	// log every: 100
}
