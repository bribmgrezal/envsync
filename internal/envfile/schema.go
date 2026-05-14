package envfile

import (
	"fmt"
	"regexp"
	"strings"
)

// SchemaField describes a single expected key in a .env file.
type SchemaField struct {
	Key      string
	Required bool
	Pattern  *regexp.Regexp // optional value validation pattern
	Default  string         // optional default value hint
}

// Schema holds a collection of field definitions used to validate env entries.
type Schema struct {
	Fields []SchemaField
}

// SchemaViolation describes a single schema validation failure.
type SchemaViolation struct {
	Key     string
	Message string
}

func (v SchemaViolation) Error() string {
	return fmt.Sprintf("schema violation for %q: %s", v.Key, v.Message)
}

// ValidateSchema checks a map of key→value pairs against the schema.
// It returns all violations found; a nil slice means the env is valid.
func ValidateSchema(s Schema, env map[string]string) []SchemaViolation {
	var violations []SchemaViolation

	for _, field := range s.Fields {
		val, present := env[field.Key]

		if !present || strings.TrimSpace(val) == "" {
			if field.Required {
				violations = append(violations, SchemaViolation{
					Key:     field.Key,
					Message: "required key is missing or empty",
				})
			}
			continue
		}

		if field.Pattern != nil && !field.Pattern.MatchString(val) {
			violations = append(violations, SchemaViolation{
				Key:     field.Key,
				Message: fmt.Sprintf("value %q does not match required pattern %s", val, field.Pattern.String()),
			})
		}
	}

	return violations
}
