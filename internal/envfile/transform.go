package envfile

import (
	"strings"
)

// TransformOptions controls how entries are transformed.
type TransformOptions struct {
	// UppercaseKeys converts all keys to uppercase.
	UppercaseKeys bool
	// TrimValues strips leading and trailing whitespace from values.
	TrimValues bool
	// PrefixKeys prepends a string to every key.
	PrefixKeys string
	// StripPrefix removes a prefix from keys that have it.
	StripPrefix string
}

// DefaultTransformOptions returns a TransformOptions with safe defaults.
func DefaultTransformOptions() TransformOptions {
	return TransformOptions{
		UppercaseKeys: false,
		TrimValues:    true,
		PrefixKeys:    "",
		StripPrefix:   "",
	}
}

// Transform applies the given options to a slice of Entry values and
// returns a new slice with the transformations applied. The original
// entries are not modified.
func Transform(entries []Entry, opts TransformOptions) []Entry {
	out := make([]Entry, 0, len(entries))
	for _, e := range entries {
		key := e.Key
		val := e.Value

		if opts.StripPrefix != "" && strings.HasPrefix(key, opts.StripPrefix) {
			key = strings.TrimPrefix(key, opts.StripPrefix)
		}
		if opts.PrefixKeys != "" {
			key = opts.PrefixKeys + key
		}
		if opts.UppercaseKeys {
			key = strings.ToUpper(key)
		}
		if opts.TrimValues {
			val = strings.TrimSpace(val)
		}

		out = append(out, Entry{
			Key:    key,
			Value:  val,
			LineNo: e.LineNo,
			Masked: e.Masked,
		})
	}
	return out
}
