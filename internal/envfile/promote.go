package envfile

import (
	"fmt"
	"strings"
)

// DefaultPromoteOptions returns a PromoteOptions with safe defaults.
func DefaultPromoteOptions() PromoteOptions {
	return PromoteOptions{
		FromPrefix: "",
		ToPrefix:   "",
		KeepOriginal: false,
		FailOnMissing: false,
	}
}

// PromoteOptions controls how keys are promoted between prefix namespaces.
type PromoteOptions struct {
	// FromPrefix is the source prefix to match (e.g. "DEV_").
	FromPrefix string

	// ToPrefix is the destination prefix to apply (e.g. "PROD_").
	ToPrefix string

	// KeepOriginal retains the original prefixed entry alongside the promoted one.
	KeepOriginal bool

	// FailOnMissing returns an error if no entries match FromPrefix.
	FailOnMissing bool
}

// Promote copies or moves entries matching FromPrefix into ToPrefix-namespaced keys.
// If FromPrefix is empty, all entries are promoted (prefix is simply prepended).
// The key stored in the result has FromPrefix stripped and ToPrefix prepended.
func Promote(entries []Entry, opts PromoteOptions) ([]Entry, error) {
	var out []Entry
	matched := 0

	for _, e := range entries {
		if opts.FromPrefix != "" && !strings.HasPrefix(e.Key, opts.FromPrefix) {
			// Not a match — keep original unconditionally.
			out = append(out, e)
			continue
		}

		matched++

		// Strip the source prefix, then apply destination prefix.
		bare := strings.TrimPrefix(e.Key, opts.FromPrefix)
		promoted := Entry{
			Key:    opts.ToPrefix + bare,
			Value:  e.Value,
			LineNo: 0, // promoted entry has no original line number
		}

		if opts.KeepOriginal {
			out = append(out, e)
		}
		out = append(out, promoted)
	}

	if opts.FailOnMissing && matched == 0 {
		return nil, fmt.Errorf("promote: no entries matched prefix %q", opts.FromPrefix)
	}

	return out, nil
}
