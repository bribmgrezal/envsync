// Package envfile provides utilities for parsing, transforming, and
// manipulating .env files.
//
// # Defaults
//
// The Defaults function fills gaps in a destination entry slice using a
// separate "defaults" slice as a fallback source.
//
// Typical usage:
//
//	base, _ := envfile.Parse("app.env")
//	fallback, _ := envfile.Parse("app.defaults.env")
//
//	opts := envfile.DefaultDefaultsOptions()
//	result, err := envfile.Defaults(base, fallback, opts)
//
// Behaviour:
//   - Keys absent from dst are appended from defaults.
//   - Keys present in dst but with an empty value are filled from defaults.
//   - Existing non-empty values are preserved unless Overwrite is true.
//   - Default entries with empty values are skipped when SkipEmpty is true
//     (the default), preventing accidental erasure of dst values.
package envfile
