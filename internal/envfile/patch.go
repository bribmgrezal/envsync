package envfile

import "fmt"

// PatchOptions controls the behaviour of the Patch function.
type PatchOptions struct {
	// AllowNew permits keys that do not exist in the base slice to be added.
	AllowNew bool
	// FailOnMissing returns an error when a patch key is absent from the base
	// and AllowNew is false, instead of silently skipping it.
	FailOnMissing bool
}

// DefaultPatchOptions returns a PatchOptions with safe defaults.
func DefaultPatchOptions() PatchOptions {
	return PatchOptions{
		AllowNew:      true,
		FailOnMissing: false,
	}
}

// Patch applies a set of key=value overrides (the patch) onto a base slice of
// Entry values and returns a new slice with the changes applied.
//
// Entries in base that are not referenced by patch are preserved unchanged.
// The relative order of base entries is maintained; new entries from patch are
// appended in the order they appear in patch.
func Patch(base, patch []Entry, opts PatchOptions) ([]Entry, error) {
	// Build an index so we can update in O(1).
	index := make(map[string]int, len(base))
	for i, e := range base {
		index[e.Key] = i
	}

	result := make([]Entry, len(base))
	copy(result, base)

	for _, p := range patch {
		if i, ok := index[p.Key]; ok {
			result[i] = Entry{
				Key:    p.Key,
				Value:  p.Value,
				LineNo: result[i].LineNo,
			}
			continue
		}

		// Key not found in base.
		if !opts.AllowNew {
			if opts.FailOnMissing {
				return nil, fmt.Errorf("patch: key %q not found in base", p.Key)
			}
			continue
		}

		result = append(result, Entry{Key: p.Key, Value: p.Value})
		index[p.Key] = len(result) - 1
	}

	return result, nil
}
