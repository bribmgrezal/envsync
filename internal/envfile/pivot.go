package envfile

import (
	"fmt"
	"strings"
)

// DefaultPivotOptions returns a PivotOptions with sensible defaults.
func DefaultPivotOptions() PivotOptions {
	return PivotOptions{
		Separator:  "_",
		KeyColumn:  "KEY",
		ValColumn:  "VALUE",
		Uppercase:  true,
	}
}

// PivotOptions controls how Pivot reshapes entries.
type PivotOptions struct {
	// Separator is used to join the namespace and field name into a new key.
	Separator string

	// KeyColumn is the entry key whose value becomes the new namespace prefix.
	KeyColumn string

	// ValColumn is the entry key whose value becomes the new entry value.
	ValColumn string

	// Uppercase forces all generated keys to upper-case.
	Uppercase bool

	// FailOnMissing returns an error when KeyColumn or ValColumn is absent in a
	// group instead of silently skipping the group.
	FailOnMissing bool
}

// Pivot reshapes a flat list of entries that represent rows of a key/value
// table into a conventional KEY=VALUE env-file shape.
//
// Each consecutive pair of entries whose keys match opts.KeyColumn and
// opts.ValColumn is collapsed into a single entry:
//
//	NAMESPACE_KEY=DB_HOST  →  DB_HOST=<value-from-ValColumn>
//
// Entries that do not match KeyColumn or ValColumn are passed through
// unchanged.
func Pivot(entries []Entry, opts PivotOptions) ([]Entry, error) {
	if opts.Separator == "" {
		opts.Separator = "_"
	}

	// Index entries by key for quick lookup.
	indexed := make(map[string][]Entry)
	var order []string
	for _, e := range entries {
		k := e.Key
		if opts.Uppercase {
			k = strings.ToUpper(k)
		}
		if _, seen := indexed[k]; !seen {
			order = append(order, k)
		}
		indexed[k] = append(indexed[k], e)
	}

	keyCol := opts.KeyColumn
	valCol := opts.ValColumn
	if opts.Uppercase {
		keyCol = strings.ToUpper(keyCol)
		valCol = strings.ToUpper(valCol)
	}

	// If neither pivot column is present, return entries unchanged.
	_, hasKey := indexed[keyCol]
	_, hasVal := indexed[valCol]
	if !hasKey && !hasVal {
		return entries, nil
	}

	if opts.FailOnMissing {
		if !hasKey {
			return nil, fmt.Errorf("pivot: key column %q not found", opts.KeyColumn)
		}
		if !hasVal {
			return nil, fmt.Errorf("pivot: value column %q not found", opts.ValColumn)
		}
	}

	var result []Entry
	for _, k := range order {
		if k == keyCol || k == valCol {
			continue
		}
		result = append(result, indexed[k]...)
	}

	// Build pivoted entries from parallel KeyColumn/ValColumn slices.
	keyEntries := indexed[keyCol]
	valEntries := indexed[valCol]
	min := len(keyEntries)
	if len(valEntries) < min {
		min = len(valEntries)
	}
	for i := 0; i < min; i++ {
		newKey := keyEntries[i].Value
		if opts.Uppercase {
			newKey = strings.ToUpper(newKey)
		}
		result = append(result, Entry{
			Key:   newKey,
			Value: valEntries[i].Value,
		})
	}

	return result, nil
}
