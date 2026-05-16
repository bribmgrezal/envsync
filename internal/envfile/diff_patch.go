package envfile

import (
	"fmt"
	"strings"
)

// PatchFromDiff represents a set of changes derived from comparing two env files.
type PatchFromDiff struct {
	Add    []Entry
	Update []Entry
	Remove []string
}

// DiffPatchOptions controls the behaviour of BuildPatchFromDiff.
type DiffPatchOptions struct {
	// IncludeRemovals causes keys present in dst but absent in src to be
	// scheduled for removal in the resulting patch.
	IncludeRemovals bool

	// IgnoreKeys is an optional list of exact key names to skip during diff.
	IgnoreKeys []string
}

// DefaultDiffPatchOptions returns a DiffPatchOptions with sensible defaults.
func DefaultDiffPatchOptions() DiffPatchOptions {
	return DiffPatchOptions{
		IncludeRemovals: false,
	}
}

// BuildPatchFromDiff compares src against dst and returns a PatchFromDiff
// describing what must change in dst to make it match src.
//
// Keys listed in opts.IgnoreKeys are skipped entirely. When
// opts.IncludeRemovals is true, keys that exist in dst but not in src are
// added to PatchFromDiff.Remove.
func BuildPatchFromDiff(src, dst []Entry, opts DiffPatchOptions) (PatchFromDiff, error) {
	if len(src) == 0 && len(dst) == 0 {
		return PatchFromDiff{}, nil
	}

	ignore := make(map[string]struct{}, len(opts.IgnoreKeys))
	for _, k := range opts.IgnoreKeys {
		key := strings.TrimSpace(k)
		if key == "" {
			return PatchFromDiff{}, fmt.Errorf("envfile: ignore key must not be blank")
		}
		ignore[key] = struct{}{}
	}

	srcMap := make(map[string]string, len(src))
	for _, e := range src {
		if _, skip := ignore[e.Key]; !skip {
			srcMap[e.Key] = e.Value
		}
	}

	dstMap := make(map[string]string, len(dst))
	for _, e := range dst {
		if _, skip := ignore[e.Key]; !skip {
			dstMap[e.Key] = e.Value
		}
	}

	var result PatchFromDiff

	for key, srcVal := range srcMap {
		dstVal, exists := dstMap[key]
		switch {
		case !exists:
			result.Add = append(result.Add, Entry{Key: key, Value: srcVal})
		case srcVal != dstVal:
			result.Update = append(result.Update, Entry{Key: key, Value: srcVal})
		}
	}

	if opts.IncludeRemovals {
		for key := range dstMap {
			if _, exists := srcMap[key]; !exists {
				result.Remove = append(result.Remove, key)
			}
		}
	}

	return result, nil
}
