// Package export provides utilities for serialising parsed .env entries
// into various output formats.
//
// Supported formats:
//
//   - dotenv  — standard KEY=VALUE lines (default)
//   - export  — shell-compatible "export KEY=VALUE" lines
//   - json    — a simple JSON object mapping keys to values
//
// Secret masking
//
// When Options.Masked is true, any entry whose Masked field is set will
// have its value replaced with "***" in the output. This allows safe
// display of environment configurations without leaking credentials.
//
// Sorting
//
// When Options.Sorted is true (the default via DefaultOptions), entries
// are emitted in ascending alphabetical key order, making diffs easier
// to read.
package export
