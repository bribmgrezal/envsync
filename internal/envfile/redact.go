package envfile

import (
	"strings"
)

// RedactOptions controls how values are redacted in a set of entries.
type RedactOptions struct {
	// Keys is an explicit list of keys whose values should be redacted.
	Keys []string
	// Patterns is a list of substring patterns; any key containing one of
	// these patterns (case-insensitive) will be redacted.
	Patterns []string
	// Placeholder is the string used to replace a redacted value.
	// Defaults to "***REDACTED***" if empty.
	Placeholder string
}

// DefaultRedactOptions returns a RedactOptions with sensible defaults.
func DefaultRedactOptions() RedactOptions {
	return RedactOptions{
		Patterns:    []string{"SECRET", "PASSWORD", "PASSWD", "TOKEN", "API_KEY", "PRIVATE", "CREDENTIAL"},
		Placeholder: "***REDACTED***",
	}
}

// Redact returns a copy of entries with sensitive values replaced by the
// configured placeholder. Original entries are never mutated.
func Redact(entries []Entry, opts RedactOptions) []Entry {
	if opts.Placeholder == "" {
		opts.Placeholder = "***REDACTED***"
	}

	keySet := make(map[string]struct{}, len(opts.Keys))
	for _, k := range opts.Keys {
		keySet[strings.ToUpper(k)] = struct{}{}
	}

	result := make([]Entry, len(entries))
	for i, e := range entries {
		copy := e
		if shouldRedact(e.Key, keySet, opts.Patterns) {
			copy.Value = opts.Placeholder
			copy.Masked = true
		}
		result[i] = copy
	}
	return result
}

func shouldRedact(key string, keySet map[string]struct{}, patterns []string) bool {
	upper := strings.ToUpper(key)
	if _, ok := keySet[upper]; ok {
		return true
	}
	for _, p := range patterns {
		if strings.Contains(upper, strings.ToUpper(p)) {
			return true
		}
	}
	return false
}
