package envfile

import "strings"

// FilterOptions controls which entries are included by Filter.
type FilterOptions struct {
	// Prefix, if non-empty, keeps only entries whose key starts with this value.
	Prefix string

	// ExcludeKeys is a set of exact key names to drop from the result.
	ExcludeKeys []string

	// SkipEmpty, when true, omits entries with empty or whitespace-only values.
	SkipEmpty bool
}

// DefaultFilterOptions returns a FilterOptions with no restrictions applied.
func DefaultFilterOptions() FilterOptions {
	return FilterOptions{}
}

// Filter returns a new slice of Entry values that satisfy all criteria
// defined in opts. The original slice is never modified.
func Filter(entries []Entry, opts FilterOptions) []Entry {
	excludeSet := make(map[string]struct{}, len(opts.ExcludeKeys))
	for _, k := range opts.ExcludeKeys {
		excludeSet[k] = struct{}{}
	}

	result := make([]Entry, 0, len(entries))
	for _, e := range entries {
		if opts.Prefix != "" && !strings.HasPrefix(e.Key, opts.Prefix) {
			continue
		}
		if _, excluded := excludeSet[e.Key]; excluded {
			continue
		}
		if opts.SkipEmpty && strings.TrimSpace(e.Value) == "" {
			continue
		}
		result = append(result, e)
	}
	return result
}
