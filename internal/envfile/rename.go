package envfile

import "fmt"

// RenameOptions controls the behaviour of the Rename function.
type RenameOptions struct {
	// Mapping is a map of old key names to new key names.
	// Keys not present in the map are left unchanged.
	Mapping map[string]string

	// SkipMissing, when true, silently ignores keys in Mapping that do not
	// exist in the source entries. When false (default) a missing key causes
	// an error to be returned.
	SkipMissing bool
}

// DefaultRenameOptions returns a RenameOptions with safe defaults.
func DefaultRenameOptions() RenameOptions {
	return RenameOptions{
		Mapping:     make(map[string]string),
		SkipMissing: true,
	}
}

// Rename returns a new slice of Entry values where keys have been renamed
// according to opts.Mapping. The order of entries is preserved; renamed
// entries keep their original position in the slice.
//
// If opts.SkipMissing is false and a key listed in opts.Mapping is not found
// in entries, Rename returns an error.
func Rename(entries []Entry, opts RenameOptions) ([]Entry, error) {
	matched := make(map[string]bool, len(opts.Mapping))

	out := make([]Entry, 0, len(entries))
	for _, e := range entries {
		if newKey, ok := opts.Mapping[e.Key]; ok {
			matched[e.Key] = true
			e.Key = newKey
		}
		out = append(out, e)
	}

	if !opts.SkipMissing {
		for oldKey := range opts.Mapping {
			if !matched[oldKey] {
				return nil, fmt.Errorf("rename: key %q not found in entries", oldKey)
			}
		}
	}

	return out, nil
}
