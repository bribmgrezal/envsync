package envfile

// DedupeOptions controls how duplicate keys are resolved during deduplication.
type DedupeOptions struct {
	// KeepFirst retains the first occurrence of a duplicate key instead of the last.
	KeepFirst bool
	// ReportDuplicates collects the keys that were removed during deduplication.
	ReportDuplicates bool
}

// DefaultDedupeOptions returns a DedupeOptions with sensible defaults:
// last occurrence wins, duplicates are not reported.
func DefaultDedupeOptions() DedupeOptions {
	return DedupeOptions{
		KeepFirst:        false,
		ReportDuplicates: false,
	}
}

// DedupeResult holds the deduplicated entries and metadata.
type DedupeResult struct {
	Entries    []Entry
	Duplicates []string // keys that had duplicates removed
}

// Dedupe removes duplicate keys from entries according to opts.
// By default the last occurrence of each key is kept, preserving
// the original relative order of the surviving entries.
func Dedupe(entries []Entry, opts DedupeOptions) DedupeResult {
	if len(entries) == 0 {
		return DedupeResult{}
	}

	type indexedEntry struct {
		entry Entry
		index int
	}

	// Track the winning index per key.
	seen := make(map[string]indexedEntry, len(entries))
	dupeSet := make(map[string]bool)

	for i, e := range entries {
		if _, exists := seen[e.Key]; exists {
			dupeSet[e.Key] = true
		}
		if opts.KeepFirst {
			if _, exists := seen[e.Key]; !exists {
				seen[e.Key] = indexedEntry{entry: e, index: i}
			}
		} else {
			seen[e.Key] = indexedEntry{entry: e, index: i}
		}
	}

	// Rebuild slice in original order using surviving indices.
	result := make([]Entry, 0, len(seen))
	visited := make(map[string]bool, len(seen))
	for _, e := range entries {
		ie := seen[e.Key]
		if !visited[e.Key] && ie.entry.Key == e.Key {
			// Only append when we reach the winning entry.
			if ie.index == indexOf(entries, e) {
				result = append(result, e)
				visited[e.Key] = true
			}
		}
	}

	var dupes []string
	if opts.ReportDuplicates {
		for k := range dupeSet {
			dupes = append(dupes, k)
		}
	}

	return DedupeResult{Entries: result, Duplicates: dupes}
}

// indexOf returns the position of e in entries by pointer-equivalent comparison.
func indexOf(entries []Entry, e Entry) int {
	for i := range entries {
		if entries[i].Key == e.Key && entries[i].Value == e.Value && entries[i].LineNo == e.LineNo {
			return i
		}
	}
	return -1
}
