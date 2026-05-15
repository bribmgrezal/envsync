package envfile

import (
	"fmt"
	"strings"
)

// FlattenOptions controls how nested key segments are joined.
type FlattenOptions struct {
	// Separator is the string used to join key segments (default: "_").
	Separator string
	// UppercaseKeys converts all resulting keys to uppercase.
	UppercaseKeys bool
	// Prefix is prepended to every flattened key (optional).
	Prefix string
}

// DefaultFlattenOptions returns sensible defaults for Flatten.
func DefaultFlattenOptions() FlattenOptions {
	return FlattenOptions{
		Separator:     "_",
		UppercaseKeys: true,
	}
}

// Flatten takes a slice of entries whose keys may contain a hierarchical
// separator (e.g. "DB.HOST" or "DB__HOST") and normalises them into a flat
// KEY=VALUE form using the configured Separator.
//
// Keys that do not contain the source separator are passed through unchanged
// (subject to UppercaseKeys and Prefix options).
//
// The source separator is auto-detected as the first occurrence of "__",
// ".", or "/" found in any key; if none is found the entries are returned
// after applying Prefix / UppercaseKeys only.
func Flatten(entries []Entry, opts FlattenOptions) []Entry {
	if opts.Separator == "" {
		opts.Separator = "_"
	}

	srcSep := detectSeparator(entries)

	out := make([]Entry, 0, len(entries))
	for _, e := range entries {
		key := e.Key
		if srcSep != "" && strings.Contains(key, srcSep) {
			parts := strings.Split(key, srcSep)
			key = strings.Join(parts, opts.Separator)
		}
		if opts.Prefix != "" {
			key = opts.Prefix + opts.Separator + key
		}
		if opts.UppercaseKeys {
			key = strings.ToUpper(key)
		}
		out = append(out, Entry{
			Key:   key,
			Value: e.Value,
		})
	}
	return out
}

// detectSeparator returns the first hierarchical separator found across all
// entry keys, preferring "__" over "." over "/".
func detectSeparator(entries []Entry) string {
	candidates := []string{"__", ".", "/"}
	for _, sep := range candidates {
		for _, e := range entries {
			if strings.Contains(e.Key, sep) {
				return sep
			}
		}
	}
	return ""
}

// FlattenMap converts a nested map[string]interface{} (e.g. parsed from JSON
// or YAML) into a flat slice of Entry values using the provided options.
func FlattenMap(m map[string]interface{}, opts FlattenOptions) []Entry {
	if opts.Separator == "" {
		opts.Separator = "_"
	}
	var out []Entry
	flattenMapRecursive(m, "", opts, &out)
	return out
}

func flattenMapRecursive(m map[string]interface{}, prefix string, opts FlattenOptions, out *[]Entry) {
	for k, v := range m {
		fullKey := k
		if prefix != "" {
			fullKey = prefix + opts.Separator + k
		}
		switch child := v.(type) {
		case map[string]interface{}:
			flattenMapRecursive(child, fullKey, opts, out)
		default:
			key := fullKey
			if opts.UppercaseKeys {
				key = strings.ToUpper(key)
			}
			*out = append(*out, Entry{
				Key:   key,
				Value: fmt.Sprintf("%v", v),
			})
		}
	}
}
