package envfile

import (
	"sort"
	"strings"
)

// SortOrder defines the ordering strategy for env entries.
type SortOrder int

const (
	// SortAlpha sorts entries alphabetically by key (A→Z).
	SortAlpha SortOrder = iota
	// SortAlphaDesc sorts entries reverse-alphabetically by key (Z→A).
	SortAlphaDesc
	// SortByGroup sorts entries by prefix group (e.g. DB_, APP_) then alphabetically within each group.
	SortByGroup
)

// DefaultSortOptions returns a SortOptions with alphabetical ordering.
func DefaultSortOptions() SortOptions {
	return SortOptions{Order: SortAlpha}
}

// SortOptions controls how entries are sorted.
type SortOptions struct {
	// Order is the sort strategy to apply.
	Order SortOrder
}

// Sort returns a new slice of entries ordered according to opts.
// The original slice is not modified.
func Sort(entries []Entry, opts SortOptions) []Entry {
	out := make([]Entry, len(entries))
	copy(out, entries)

	switch opts.Order {
	case SortAlphaDesc:
		sort.SliceStable(out, func(i, j int) bool {
			return out[i].Key > out[j].Key
		})
	case SortByGroup:
		sort.SliceStable(out, func(i, j int) bool {
			gi := groupPrefix(out[i].Key)
			gj := groupPrefix(out[j].Key)
			if gi != gj {
				return gi < gj
			}
			return out[i].Key < out[j].Key
		})
	default: // SortAlpha
		sort.SliceStable(out, func(i, j int) bool {
			return out[i].Key < out[j].Key
		})
	}

	return out
}

// groupPrefix returns the portion of key up to and including the first '_',
// or the whole key if no underscore is present.
func groupPrefix(key string) string {
	if idx := strings.Index(key, "_"); idx >= 0 {
		return key[:idx+1]
	}
	return key
}
