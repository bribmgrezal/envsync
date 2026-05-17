package envfile

// CoalesceOptions controls the behaviour of Coalesce.
type CoalesceOptions struct {
	// PreferFirst returns the first non-empty value found across sources
	// instead of the last. Defaults to false (last wins).
	PreferFirst bool

	// SkipEmpty treats empty-string values as absent so that a later
	// (or earlier, when PreferFirst is true) source can fill the gap.
	SkipEmpty bool
}

// DefaultCoalesceOptions returns a CoalesceOptions with sensible defaults.
func DefaultCoalesceOptions() CoalesceOptions {
	return CoalesceOptions{
		PreferFirst: false,
		SkipEmpty:   true,
	}
}

// Coalesce merges multiple slices of Entry into a single slice.
// Keys that appear in more than one source are resolved according to opts:
//   - When PreferFirst is false (default) the last source that carries a
//     non-empty value wins.
//   - When PreferFirst is true the first source that carries a non-empty
//     value wins and subsequent sources are ignored for that key.
//
// Source order is preserved for keys that appear only once. For keys that
// appear in multiple sources the position of the winning entry is used.
func Coalesce(sources [][]Entry, opts CoalesceOptions) []Entry {
	type slot struct {
		entry Entry
		order int
	}

	seen := make(map[string]slot)
	counter := 0

	for _, src := range sources {
		for _, e := range src {
			if opts.SkipEmpty && e.Value == "" {
				continue
			}
			existing, exists := seen[e.Key]
			if !exists {
				seen[e.Key] = slot{entry: e, order: counter}
				counter++
				continue
			}
			if opts.PreferFirst {
				// keep existing — first wins
				_ = existing
				continue
			}
			// last wins: update value but keep original insertion order
			seen[e.Key] = slot{entry: e, order: existing.order}
		}
	}

	// Reconstruct in insertion order.
	result := make([]Entry, len(seen))
	for _, s := range seen {
		result[s.order] = s.entry
	}
	return result
}
