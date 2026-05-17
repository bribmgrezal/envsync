package envfile

import "strings"

// MaskOptions controls how values are masked during output.
type MaskOptions struct {
	// Placeholder replaces sensitive values. Defaults to "****".
	Placeholder string
	// Patterns is a list of key substrings treated as sensitive.
	// If empty, a built-in list is used.
	Patterns []string
	// ExplicitKeys forces masking for these keys regardless of pattern matching.
	ExplicitKeys []string
}

// DefaultMaskOptions returns sensible defaults for MaskOptions.
func DefaultMaskOptions() MaskOptions {
	return MaskOptions{
		Placeholder: "****",
		Patterns:    []string{"SECRET", "PASSWORD", "PASSWD", "TOKEN", "API_KEY", "PRIVATE", "CREDENTIAL"},
	}
}

// MaskValues returns a copy of entries where sensitive values are replaced
// with the configured placeholder. The original slice is not mutated.
func MaskValues(entries []Entry, opts MaskOptions) []Entry {
	if opts.Placeholder == "" {
		opts.Placeholder = "****"
	}

	explicit := make(map[string]struct{}, len(opts.ExplicitKeys))
	for _, k := range opts.ExplicitKeys {
		explicit[strings.ToUpper(k)] = struct{}{}
	}

	out := make([]Entry, len(entries))
	for i, e := range entries {
		copy := e
		if isMaskSensitive(e.Key, opts.Patterns, explicit) {
			copy.Value = opts.Placeholder
			copy.Masked = true
		}
		out[i] = copy
	}
	return out
}

func isMaskSensitive(key string, patterns []string, explicit map[string]struct{}) bool {
	upper := strings.ToUpper(key)
	if _, ok := explicit[upper]; ok {
		return true
	}
	for _, p := range patterns {
		if strings.Contains(upper, p) {
			return true
		}
	}
	return false
}
