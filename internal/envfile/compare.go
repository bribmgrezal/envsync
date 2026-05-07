package envfile

import "sort"

// DiffEntry represents a single difference between two env files.
type DiffEntry struct {
	Key      string
	BaseVal  string
	OtherVal string
	Status   DiffStatus
}

// DiffStatus indicates the type of difference found.
type DiffStatus int

const (
	DiffAdded   DiffStatus = iota // Key exists in other but not in base
	DiffRemoved                   // Key exists in base but not in other
	DiffChanged                   // Key exists in both but values differ
	DiffSame                      // Key exists in both with identical values
)

// String returns a human-readable label for the diff status.
func (s DiffStatus) String() string {
	switch s {
	case DiffAdded:
		return "added"
	case DiffRemoved:
		return "removed"
	case DiffChanged:
		return "changed"
	default:
		return "same"
	}
}

// CompareEntries diffs two slices of Entry by key, returning a sorted list
// of DiffEntry values describing additions, removals, changes, and matches.
func CompareEntries(base, other []Entry) []DiffEntry {
	baseMap := ToMap(base)
	otherMap := ToMap(other)

	seen := make(map[string]bool)
	var results []DiffEntry

	for k, bv := range baseMap {
		seen[k] = true
		if ov, ok := otherMap[k]; ok {
			if bv == ov {
				results = append(results, DiffEntry{Key: k, BaseVal: bv, OtherVal: ov, Status: DiffSame})
			} else {
				results = append(results, DiffEntry{Key: k, BaseVal: bv, OtherVal: ov, Status: DiffChanged})
			}
		} else {
			results = append(results, DiffEntry{Key: k, BaseVal: bv, Status: DiffRemoved})
		}
	}

	for k, ov := range otherMap {
		if !seen[k] {
			results = append(results, DiffEntry{Key: k, OtherVal: ov, Status: DiffAdded})
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Key < results[j].Key
	})

	return results
}
