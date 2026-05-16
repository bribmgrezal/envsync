package envfile

import (
	"fmt"
	"strconv"
	"strings"
)

// CastType represents the target type for a value cast.
type CastType string

const (
	CastString CastType = "string"
	CastInt    CastType = "int"
	CastFloat  CastType = "float"
	CastBool   CastType = "bool"
)

// CastRule maps a key to a desired CastType.
type CastRule struct {
	Key  string
	Type CastType
}

// DefaultCastOptions returns options with no rules and strict mode enabled.
func DefaultCastOptions() CastOptions {
	return CastOptions{
		Rules:  nil,
		Strict: true,
	}
}

// CastOptions controls how Cast behaves.
type CastOptions struct {
	// Rules defines per-key type coercions.
	Rules []CastRule
	// Strict causes Cast to return an error if a value cannot be coerced.
	// When false, the original value is kept on failure.
	Strict bool
}

// CastResult records the outcome of casting a single entry.
type CastResult struct {
	Key      string
	Original string
	Casted   string
	Type     CastType
	Err      error
}

// Cast normalises env entry values to canonical string representations of the
// target type (e.g. "true" instead of "TRUE", "42" instead of "042").
// It returns the modified entries and a slice of per-key results.
func Cast(entries []Entry, opts CastOptions) ([]Entry, []CastResult, error) {
	ruleMap := make(map[string]CastType, len(opts.Rules))
	for _, r := range opts.Rules {
		ruleMap[r.Key] = r.Type
	}

	out := make([]Entry, len(entries))
	copy(out, entries)

	var results []CastResult

	for i, e := range out {
		typ, ok := ruleMap[e.Key]
		if !ok {
			continue
		}
		casted, err := coerce(e.Value, typ)
		results = append(results, CastResult{
			Key:      e.Key,
			Original: e.Value,
			Casted:   casted,
			Type:     typ,
			Err:      err,
		})
		if err != nil {
			if opts.Strict {
				return nil, results, fmt.Errorf("cast %q to %s: %w", e.Key, typ, err)
			}
			continue
		}
		out[i].Value = casted
	}
	return out, results, nil
}

func coerce(val string, typ CastType) (string, error) {
	switch typ {
	case CastString:
		return val, nil
	case CastInt:
		n, err := strconv.ParseInt(strings.TrimSpace(val), 0, 64)
		if err != nil {
			return val, fmt.Errorf("cannot parse %q as int: %w", val, err)
		}
		return strconv.FormatInt(n, 10), nil
	case CastFloat:
		f, err := strconv.ParseFloat(strings.TrimSpace(val), 64)
		if err != nil {
			return val, fmt.Errorf("cannot parse %q as float: %w", val, err)
		}
		return strconv.FormatFloat(f, 'f', -1, 64), nil
	case CastBool:
		b, err := strconv.ParseBool(strings.TrimSpace(val))
		if err != nil {
			return val, fmt.Errorf("cannot parse %q as bool: %w", val, err)
		}
		return strconv.FormatBool(b), nil
	default:
		return val, fmt.Errorf("unknown cast type %q", typ)
	}
}
