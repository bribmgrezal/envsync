package envfile

import "strings"

// TrimOptions controls which trimming operations are applied to entries.
type TrimOptions struct {
	// TrimKeys removes leading and trailing whitespace from all keys.
	TrimKeys bool
	// TrimValues removes leading and trailing whitespace from all values.
	TrimValues bool
	// RemoveQuotes strips surrounding single or double quotes from values.
	RemoveQuotes bool
	// NormalizeKeys converts all keys to uppercase and replaces hyphens with underscores.
	NormalizeKeys bool
}

// DefaultTrimOptions returns a TrimOptions with sensible defaults enabled.
func DefaultTrimOptions() TrimOptions {
	return TrimOptions{
		TrimKeys:     true,
		TrimValues:   true,
		RemoveQuotes: false,
		NormalizeKeys: false,
	}
}

// Trim applies the configured trimming and normalisation operations to a
// slice of Entry values. The original slice is not mutated; a new slice
// is returned.
func Trim(entries []Entry, opts TrimOptions) []Entry {
	out := make([]Entry, 0, len(entries))
	for _, e := range entries {
		if opts.TrimKeys {
			e.Key = strings.TrimSpace(e.Key)
		}
		if opts.TrimValues {
			e.Value = strings.TrimSpace(e.Value)
		}
		if opts.RemoveQuotes {
			e.Value = stripQuotes(e.Value)
		}
		if opts.NormalizeKeys {
			e.Key = strings.ToUpper(strings.ReplaceAll(e.Key, "-", "_"))
		}
		out = append(out, e)
	}
	return out
}

// stripQuotes removes a matching pair of surrounding single or double quotes
// from s. If s is not quoted, it is returned unchanged.
func stripQuotes(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') ||
			(s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
