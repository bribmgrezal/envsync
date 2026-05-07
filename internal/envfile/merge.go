package envfile

// MergeOptions controls how two sets of env entries are merged.
type MergeOptions struct {
	// Overwrite determines whether values from src overwrite existing keys in dst.
	Overwrite bool
	// SkipEmpty skips entries from src whose value is empty.
	SkipEmpty bool
}

// DefaultMergeOptions returns sensible defaults for merging.
func DefaultMergeOptions() MergeOptions {
	return MergeOptions{
		Overwrite: true,
		SkipEmpty: false,
	}
}

// Merge combines src entries into dst entries according to opts.
// dst entries are treated as the base; src entries are layered on top.
// The returned slice preserves the original order of dst, appending new keys
// from src at the end.
func Merge(dst, src []Entry, opts MergeOptions) []Entry {
	// Build an index of dst by key for fast lookup.
	index := make(map[string]int, len(dst))
	for i, e := range dst {
		index[e.Key] = i
	}

	result := make([]Entry, len(dst))
	copy(result, dst)

	for _, s := range src {
		if opts.SkipEmpty && s.Value == "" {
			continue
		}
		if i, exists := index[s.Key]; exists {
			if opts.Overwrite {
				result[i].Value = s.Value
			}
		} else {
			// New key from src — append and track its position.
			index[s.Key] = len(result)
			result = append(result, s)
		}
	}

	return result
}
