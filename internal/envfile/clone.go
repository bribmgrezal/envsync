package envfile

// CloneOptions controls how entries are cloned/duplicated.
type CloneOptions struct {
	// KeyMap maps source keys to one or more destination keys.
	// Each source key can produce multiple cloned entries.
	KeyMap map[string][]string

	// SkipMissing, when true, silently ignores source keys not found
	// in the input entries. When false (default), an error is returned.
	SkipMissing bool

	// KeepOriginal, when true, retains the original entry alongside
	// the cloned entries. Defaults to true.
	KeepOriginal bool
}

// DefaultCloneOptions returns a CloneOptions with safe defaults.
func DefaultCloneOptions() CloneOptions {
	return CloneOptions{
		SkipMissing:  false,
		KeepOriginal: true,
	}
}

// Clone duplicates entries according to the provided KeyMap.
// Each entry whose key appears in KeyMap is copied under each of the
// mapped destination names. The relative order of entries is preserved;
// cloned entries are inserted immediately after their source.
func Clone(entries []Entry, opts CloneOptions) ([]Entry, error) {
	if opts.KeyMap == nil {
		opts.KeyMap = map[string][]string{}
	}

	// Build a lookup of which source keys exist.
	existing := make(map[string]bool, len(entries))
	for _, e := range entries {
		existing[e.Key] = true
	}

	if !opts.SkipMissing {
		for src := range opts.KeyMap {
			if !existing[src] {
				return nil, fmt.Errorf("clone: source key %q not found in entries", src)
			}
		}
	}

	out := make([]Entry, 0, len(entries))
	for _, e := range entries {
		if opts.KeepOriginal {
			out = append(out, e)
		}
		destKeys, mapped := opts.KeyMap[e.Key]
		if !mapped {
			if !opts.KeepOriginal {
				out = append(out, e)
			}
			continue
		}
		for _, dk := range destKeys {
			cloned := e
			cloned.Key = dk
			cloned.LineNo = 0
			out = append(out, cloned)
		}
	}
	return out, nil
}
