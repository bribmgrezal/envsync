package envfile

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// interpolatePattern matches ${VAR} and $VAR style references.
var interpolatePattern = regexp.MustCompile(`\$\{([A-Za-z_][A-Za-z0-9_]*)\}|\$([A-Za-z_][A-Za-z0-9_]*)`)

// InterpolateOptions controls how variable interpolation is performed.
type InterpolateOptions struct {
	// UseOS allows falling back to OS environment variables when a key is
	// not found in the current entry set.
	UseOS bool

	// FailOnMissing returns an error when a referenced variable cannot be
	// resolved instead of leaving the placeholder intact.
	FailOnMissing bool
}

// DefaultInterpolateOptions returns sensible defaults for interpolation.
func DefaultInterpolateOptions() InterpolateOptions {
	return InterpolateOptions{
		UseOS:         true,
		FailOnMissing: false,
	}
}

// Interpolate resolves variable references within entry values.
// References of the form ${VAR} or $VAR are replaced with the corresponding
// value from the entries map or, when UseOS is enabled, from the OS environment.
func Interpolate(entries []Entry, opts InterpolateOptions) ([]Entry, error) {
	lookup := make(map[string]string, len(entries))
	for _, e := range entries {
		lookup[e.Key] = e.Value
	}

	result := make([]Entry, len(entries))
	for i, e := range entries {
		expanded, err := expand(e.Value, lookup, opts)
		if err != nil {
			return nil, fmt.Errorf("interpolate %q: %w", e.Key, err)
		}
		e.Value = expanded
		result[i] = e
	}
	return result, nil
}

func expand(value string, lookup map[string]string, opts InterpolateOptions) (string, error) {
	var expandErr error
	result := interpolatePattern.ReplaceAllStringFunc(value, func(match string) string {
		if expandErr != nil {
			return match
		}
		key := strings.TrimPrefix(strings.Trim(match, "${}"), "$")
		if v, ok := lookup[key]; ok {
			return v
		}
		if opts.UseOS {
			if v, ok := os.LookupEnv(key); ok {
				return v
			}
		}
		if opts.FailOnMissing {
			expandErr = fmt.Errorf("undefined variable %q", key)
			return match
		}
		return match
	})
	if expandErr != nil {
		return "", expandErr
	}
	return result, nil
}
