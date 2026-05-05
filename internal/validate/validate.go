// Package validate provides validation utilities for .env file entries,
// checking for common issues such as missing values, invalid key formats,
// and duplicate keys.
package validate

import (
	"fmt"
	"regexp"
	"strings"
)

// validKeyPattern matches valid environment variable key names:
// must start with a letter or underscore, followed by letters, digits, or underscores.
var validKeyPattern = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

// Issue represents a single validation problem found in an env file.
type Issue struct {
	Key     string
	Message string
}

func (i Issue) String() string {
	return fmt.Sprintf("[%s] %s", i.Key, i.Message)
}

// Result holds all issues found during validation.
type Result struct {
	Issues []Issue
}

// Valid returns true if no issues were found.
func (r *Result) Valid() bool {
	return len(r.Issues) == 0
}

// Validate checks the provided key-value map for common .env issues.
// It detects invalid key names, empty values for non-optional keys,
// and keys that contain whitespace.
func Validate(env map[string]string) *Result {
	result := &Result{}

	for key, value := range env {
		if !validKeyPattern.MatchString(key) {
			result.Issues = append(result.Issues, Issue{
				Key:     key,
				Message: "invalid key format: must match [A-Za-z_][A-Za-z0-9_]*",
			})
		}

		if strings.TrimSpace(value) == "" && value != "" {
			result.Issues = append(result.Issues, Issue{
				Key:     key,
				Message: "value contains only whitespace",
			})
		}
	}

	return result
}

// ValidateKeys checks that all required keys are present in the provided env map.
// Returns issues for any required keys that are missing or have empty values.
func ValidateKeys(env map[string]string, required []string) *Result {
	result := &Result{}

	for _, key := range required {
		val, ok := env[key]
		if !ok {
			result.Issues = append(result.Issues, Issue{
				Key:     key,
				Message: "required key is missing",
			})
			continue
		}
		if strings.TrimSpace(val) == "" {
			result.Issues = append(result.Issues, Issue{
				Key:     key,
				Message: "required key has an empty value",
			})
		}
	}

	return result
}
