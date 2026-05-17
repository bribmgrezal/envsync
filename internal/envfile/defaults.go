package envfile

import "fmt"

// DefaultsOptions controls the behaviour of the Defaults function.
type DefaultsOptions struct {
	// Overwrite replaces existing non-empty values with the default.
	// When false (the default) only missing or empty keys are filled.
	Overwrite bool

	// SkipEmpty ignores default entries whose value is empty.
	SkipEmpty bool
}

// DefaultDefaultsOptions returns a DefaultsOptions with sensible defaults.
func DefaultDefaultsOptions() DefaultsOptions {
	return DefaultsOptions{
		Overwrite: false,
		SkipEmpty: true,
	}
}

// Defaults fills missing or empty keys in dst using values from defaults.
// It returns a new slice; neither dst nor defaults is mutated.
func Defaults(dst []Entry, defaults []Entry, opts DefaultsOptions) ([]Entry, error) {
	if defaults == nil {
		return nil, fmt.Errorf("defaults: defaults slice must not be nil")
	}

	// Build a lookup of existing dst values.
	dstIndex := make(map[string]int, len(dst))
	out := make([]Entry, len(dst))
	copy(out, dst)
	for i, e := range out {
		dstIndex[e.Key] = i
	}

	for _, def := range defaults {
		if opts.SkipEmpty && def.Value == "" {
			continue
		}

		idx, exists := dstIndex[def.Key]
		if !exists {
			// Key absent — append it.
			out = append(out, Entry{Key: def.Key, Value: def.Value})
			dstIndex[def.Key] = len(out) - 1
			continue
		}

		if opts.Overwrite || out[idx].Value == "" {
			out[idx].Value = def.Value
		}
	}

	return out, nil
}
