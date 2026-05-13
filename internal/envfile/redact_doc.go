// Package envfile provides utilities for parsing, manipulating, and writing
// .env files.
//
// # Redact
//
// The Redact function replaces sensitive values in a slice of Entry with a
// configurable placeholder string. It is non-destructive: the original entries
// are never modified.
//
// Sensitivity is determined by two independent mechanisms:
//
//  1. Explicit key names supplied via RedactOptions.Keys.
//  2. Substring pattern matching against key names via RedactOptions.Patterns.
//     Matching is case-insensitive, so "password" matches "DB_PASSWORD".
//
// DefaultRedactOptions returns a ready-to-use configuration that covers the
// most common secret key naming conventions (SECRET, PASSWORD, TOKEN, API_KEY,
// PRIVATE, CREDENTIAL).
//
// Example:
//
//	entries, _ := envfile.Parse(".env")
//	redacted := envfile.Redact(entries, envfile.DefaultRedactOptions())
//
// When an entry is redacted its Masked field is set to true, which downstream
// formatters (export, report) can use to apply additional visual treatment.
package envfile
