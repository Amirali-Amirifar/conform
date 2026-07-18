# conform

**Parse, don't validate - as a reusable pattern for Go.**

`conform` is a small library for expressing domain invariants as types. 
It brings the Parse, Don't Validate pattern to Go without requiring pages of 
constructors and custom unmarshalers.

Go's type system can't say *"an int between 0 and 100"* or *"a non-empty hostname"*.
Those facts end up in doc comments, validator tags, and `if` statements scattered
through the codebase, While the type `int` promises nothing.
Every function downstream must re-check the value or trust it,
and can't tell from the type whether validation ever happened.
That's the [shotgun parsing](https://lexi-lambda.github.io/blog/2019/11/05/parse-don-t-validate/) anti-pattern.

The known remedy is a newtype with a constructor and hand-written Unmarshals,
which almost nobody writes, because it's a page of boilerplate per type.
`conform` is that pattern with the boilerplate factored out: declare rules once,
next to the field, and construction, decoding, and error reporting all enforce them.

```go
type ServerConfig struct {
    Workers  conform.Int[int] `json:"workers"`
    LogEvery conform.Int[int] `json:"log_every"`
}

func NewServerConfig() ServerConfig {
    return ServerConfig{
        Workers:  conform.NewInt(conform.Min(1), conform.Max(64)),
        LogEvery: conform.NewInt(conform.In(1, 10, 100, 1000)),
    }
}

cfg := NewServerConfig()
if err := json.Unmarshal(data, &cfg); err != nil {
    log.Fatal(err) // e.g. "got 200 but need at most 64"
}
workers := cfg.Workers.Value() // plain int, checked against 1..64
```

## How it works

A conform type is a box holding rules and, after parsing, a value:

- **The only way in is `Parse`.** It runs the rules and returns an error instead
of a value when they fail. Deserialization route through it, so for every field
present in the input, decoding is validating.
- **Misuse fails loudly.** A zero-value box rejects everything;
reading an unparsed box panics; marshaling one errors.
- **Rules are compiled Go**, not strings in tags: `func(v T) error`, written
next to the field, composed by listing them.

To name an invariant Go has no type for, wrap it:

```go
// 0..100 - no Go integer type can promise this.
type Percent struct{ conform.Int[int] }

func NewPercent() Percent {
    return Percent{conform.NewInt(conform.Min(0), conform.Max(100))}
}

func setVolume(p Percent) // the signature states the requirement
```

Embedding keeps the unmarshalers, so a `Percent` field still validates itself during decoding.

## Status

This project is in its early stages.

The API should be considered experimental and may change without notice.
Only a small set of constrained types are currently implemented, and the
library has not yet been used in production.

```sh
go get github.com/amirali-amirifar/conform
```

## License

[MIT](LICENSE)
