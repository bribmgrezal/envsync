package envfile

import (
	"fmt"
	"strconv"
	"strings"
)

// TypeRule describes the expected type for a single key.
type TypeRule struct {
	Key      string
	Type     string // "int", "float", "bool", "string", "url", "email"
	Required bool
}

// TypeViolation records a key that failed type checking.
type TypeViolation struct {
	Key      string
	Value    string
	Expected string
	Reason   string
}

// DefaultTypeCheckOptions returns a TypeCheckOptions with sensible defaults.
func DefaultTypeCheckOptions() TypeCheckOptions {
	return TypeCheckOptions{StopOnFirst: false}
}

// TypeCheckOptions controls the behaviour of TypeCheck.
type TypeCheckOptions struct {
	StopOnFirst bool // stop after the first violation
}

// TypeCheck validates each Entry whose key appears in rules against the
// declared type. Entries whose keys are not listed in rules are ignored.
func TypeCheck(entries []Entry, rules []TypeRule, opts TypeCheckOptions) ([]TypeViolation, error) {
	ruleMap := make(map[string]TypeRule, len(rules))
	for _, r := range rules {
		ruleMap[strings.ToUpper(r.Key)] = r
	}

	present := make(map[string]bool, len(entries))
	for _, e := range entries {
		present[strings.ToUpper(e.Key)] = true
	}

	var violations []TypeViolation

	// Check required keys that are absent.
	for _, r := range rules {
		if r.Required && !present[strings.ToUpper(r.Key)] {
			v := TypeViolation{Key: r.Key, Expected: r.Type, Reason: "required key is missing"}
			violations = append(violations, v)
			if opts.StopOnFirst {
				return violations, nil
			}
		}
	}

	// Check type constraints for present entries.
	for _, e := range entries {
		r, ok := ruleMap[strings.ToUpper(e.Key)]
		if !ok {
			continue
		}
		if reason := checkType(e.Value, r.Type); reason != "" {
			v := TypeViolation{Key: e.Key, Value: e.Value, Expected: r.Type, Reason: reason}
			violations = append(violations, v)
			if opts.StopOnFirst {
				return violations, nil
			}
		}
	}

	return violations, nil
}

func checkType(value, typ string) string {
	switch strings.ToLower(typ) {
	case "int":
		if _, err := strconv.ParseInt(value, 10, 64); err != nil {
			return fmt.Sprintf("expected int, got %q", value)
		}
	case "float":
		if _, err := strconv.ParseFloat(value, 64); err != nil {
			return fmt.Sprintf("expected float, got %q", value)
		}
	case "bool":
		lower := strings.ToLower(value)
		if lower != "true" && lower != "false" && lower != "1" && lower != "0" {
			return fmt.Sprintf("expected bool (true/false/1/0), got %q", value)
		}
	case "url":
		if !strings.HasPrefix(value, "http://") && !strings.HasPrefix(value, "https://") {
			return fmt.Sprintf("expected URL starting with http:// or https://, got %q", value)
		}
	case "email":
		if !strings.Contains(value, "@") || strings.HasPrefix(value, "@") || strings.HasSuffix(value, "@") {
			return fmt.Sprintf("expected email address, got %q", value)
		}
	case "string":
		// any value is valid
	}
	return ""
}
