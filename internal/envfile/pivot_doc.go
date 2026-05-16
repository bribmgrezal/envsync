// Package envfile provides utilities for parsing, transforming, and writing
// .env files.
//
// # Pivot
//
// Pivot reshapes a flat list of entries that encode a key/value table using
// two dedicated columns (KeyColumn and ValColumn) into a conventional
// KEY=VALUE env-file layout.
//
// This is useful when consuming tabular data exported from a secrets manager
// or configuration store that serialises rows as repeated KEY/VALUE pairs:
//
//	KEY=DB_HOST
//	VALUE=localhost
//	KEY=DB_PORT
//	VALUE=5432
//
// After pivoting the result is:
//
//	DB_HOST=localhost
//	DB_PORT=5432
//
// # Options
//
//   - Separator   – joined between namespace and field (default "_").
//   - KeyColumn   – entry key that holds the new key name (default "KEY").
//   - ValColumn   – entry key that holds the new value (default "VALUE").
//   - Uppercase   – normalise generated keys to upper-case (default true).
//   - FailOnMissing – return an error when a pivot column is absent instead of
//     silently passing entries through.
package envfile
