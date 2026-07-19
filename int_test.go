package conform_test

import (
	"encoding/json"
	"testing"

	"github.com/amirali-amirifar/conform"
	"github.com/amirali-amirifar/conform/predicate"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

type S struct {
	Val conform.Int[int8] `json:"val" yaml:"val"`
}

// newVal builds the spec used across most tests: 1 <= v <= 100.
func newVal() conform.Int[int8] {
	return conform.NewInt(predicate.NewAnd(conform.Min[int8](1), conform.Max[int8](100)))
}

func TestUnmarshal(t *testing.T) {
	tests := []struct {
		name      string
		doc       string
		unmarshal func([]byte, any) error
		wantErr   bool
		want      int8
	}{
		{"JSON valid", `{"val":89}`, json.Unmarshal, false, 89},
		{"JSON out of range", `{"val":120}`, json.Unmarshal, true, 0},
		{"YAML valid", "val: 89", yaml.Unmarshal, false, 89},
		{"YAML out of range", "val: 120", yaml.Unmarshal, true, 0},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := require.New(t)
			cfg := S{Val: newVal()}
			err := tc.unmarshal([]byte(tc.doc), &cfg)
			if tc.wantErr {
				r.Error(err)
				return
			}
			r.NoError(err)
			r.EqualValues(tc.want, cfg.Val.Value())
		})
	}
}

func TestMarshal(t *testing.T) {
	r := require.New(t)
	for _, marshal := range []func(any) ([]byte, error){json.Marshal, yaml.Marshal} {
		cfg := S{Val: newVal()}
		r.NoError(json.Unmarshal([]byte(`{"val":89}`), &cfg))
		out, err := marshal(cfg)
		r.NoError(err)
		r.Contains(string(out), "89") // value survives round trip

		_, err = marshal(S{Val: newVal()}) // never parsed
		r.Error(err)
	}
}

// A zero-value Int (never built via NewInt) must reject unmarshal — the
// proof-type invariant: only constructed specs can hold a value.
func TestZeroValueRejected(t *testing.T) {
	r := require.New(t)
	var cfg S
	r.Error(json.Unmarshal([]byte(`{"val":42}`), &cfg))
	r.False(cfg.Val.IsValid())
}

// Every failing rule is reported, not just the first.
func TestCollectAllRules(t *testing.T) {
	r := require.New(t)
	cfg := struct {
		Val conform.Int[int8] `json:"val"`
	}{Val: conform.NewInt(predicate.NewAnd(conform.Min[int8](10), conform.In[int8](1, 2, 3)))}

	err := json.Unmarshal([]byte(`{"val":5}`), &cfg)
	r.Error(err)
	r.ErrorContains(err, "at least 10")
	r.ErrorContains(err, "one of")
}

// A configured spec with a nil root accepts any value.
func TestNoRules(t *testing.T) {
	r := require.New(t)
	cfg := struct {
		Val conform.Int[int8] `json:"val"`
	}{Val: conform.NewInt[int8](nil)}

	r.NoError(json.Unmarshal([]byte(`{"val":42}`), &cfg))
	r.EqualValues(42, cfg.Val.Value())
}
