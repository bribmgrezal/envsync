package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"regexp"

	"github.com/user/envsync/internal/envfile"
)

// runSchema validates a .env file against a JSON schema definition file.
// Schema JSON format: [{"key":"PORT","required":true,"pattern":"^\\d+$"},...].
// Exits with a non-zero status when any violation is found.
func runSchema(args []string) error {
	fs := flag.NewFlagSet("schema", flag.ContinueOnError)
	schemaPath := fs.String("schema", "", "path to JSON schema file (required)")
	envPath := fs.String("env", ".env", "path to .env file")

	if err := fs.Parse(args); err != nil {
		return err
	}
	if *schemaPath == "" {
		return fmt.Errorf("--schema is required")
	}

	entries, err := envfile.Parse(*envPath)
	if err != nil {
		return fmt.Errorf("parse env: %w", err)
	}
	env := envfile.ToMap(entries)

	schema, err := loadSchema(*schemaPath)
	if err != nil {
		return fmt.Errorf("load schema: %w", err)
	}

	violations := envfile.ValidateSchema(schema, env)
	if len(violations) == 0 {
		fmt.Println("schema validation passed — no violations found")
		return nil
	}

	fmt.Fprintf(os.Stderr, "schema validation failed (%d violation(s)):\n", len(violations))
	for _, v := range violations {
		fmt.Fprintf(os.Stderr, "  [%s] %s\n", v.Key, v.Message)
	}
	return fmt.Errorf("schema validation failed with %d violation(s)", len(violations))
}

type rawSchemaField struct {
	Key      string `json:"key"`
	Required bool   `json:"required"`
	Pattern  string `json:"pattern"`
	Default  string `json:"default"`
}

// loadSchema reads and parses a JSON schema file from the given path.
// Each field entry may specify a key name, whether it is required, an optional
// regex pattern the value must match, and an optional default value.
func loadSchema(path string) (envfile.Schema, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return envfile.Schema{}, fmt.Errorf("read schema file %q: %w", path, err)
	}
	var raw []rawSchemaField
	if err := json.Unmarshal(data, &raw); err != nil {
		return envfile.Schema{}, fmt.Errorf("invalid schema JSON in %q: %w", path, err)
	}
	if len(raw) == 0 {
		return envfile.Schema{}, fmt.Errorf("schema file %q contains no fields", path)
	}
	var s envfile.Schema
	for _, r := range raw {
		if r.Key == "" {
			return envfile.Schema{}, fmt.Errorf("schema entry is missing a \"key\" field")
		}
		f := envfile.SchemaField{
			Key:      r.Key,
			Required: r.Required,
			Default:  r.Default,
		}
		if r.Pattern != "" {
			pat, err := regexp.Compile(r.Pattern)
			if err != nil {
				return envfile.Schema{}, fmt.Errorf("invalid pattern for key %q: %w", r.Key, err)
			}
			f.Pattern = pat
		}
		s.Fields = append(s.Fields, f)
	}
	return s, nil
}
