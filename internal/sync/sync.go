// Package sync provides functionality to apply diff results
// from one environment file to another, optionally prompting
// for missing keys or overwriting changed values.
package sync

import (
	"fmt"
	"strings"

	"github.com/user/envsync/internal/diff"
)

// Options configures the behaviour of the Sync operation.
type Options struct {
	// AddMissing controls whether keys missing in the target are added.
	AddMissing bool
	// UpdateChanged controls whether keys with changed values are updated.
	UpdateChanged bool
	// RemoveExtra controls whether keys present only in the target are removed.
	RemoveExtra bool
}

// DefaultOptions returns a sensible default: add missing, update changed,
// leave extra keys in place.
func DefaultOptions() Options {
	return Options{
		AddMissing:    true,
		UpdateChanged: true,
		RemoveExtra:   false,
	}
}

// Apply takes a slice of diff.Result entries and a target env map,
// applies changes according to opts, and returns the updated map together
// with a human-readable summary of what was changed.
func Apply(results []diff.Result, target map[string]string, opts Options) (map[string]string, string) {
	out := make(map[string]string, len(target))
	for k, v := range target {
		out[k] = v
	}

	var lines []string

	for _, r := range results {
		switch r.Status {
		case diff.Missing:
			if opts.AddMissing {
				out[r.Key] = r.SourceValue
				lines = append(lines, fmt.Sprintf("+ %s (added)", r.Key))
			}
		case diff.Changed:
			if opts.UpdateChanged {
				out[r.Key] = r.SourceValue
				lines = append(lines, fmt.Sprintf("~ %s (updated)", r.Key))
			}
		case diff.Extra:
			if opts.RemoveExtra {
				delete(out, r.Key)
				lines = append(lines, fmt.Sprintf("- %s (removed)", r.Key))
			}
		}
	}

	if len(lines) == 0 {
		return out, "no changes applied"
	}
	return out, strings.Join(lines, "\n")
}
