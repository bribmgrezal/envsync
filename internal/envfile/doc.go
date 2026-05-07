// Package envfile provides utilities for reading, writing, and manipulating
// .env files.
//
// The package supports:
//
//   - Parsing .env files into a slice of Entry values (Parse)
//   - Converting between Entry slices and maps (ToMap, FromMap)
//   - Filtering entries by prefix, exclusion list, or empty-value policy (Filter)
//   - Merging two sets of entries with configurable overwrite behaviour (Merge)
//   - Transforming entries by uppercasing keys, trimming values, or
//     adding/removing key prefixes (Transform)
//
// Entry order is preserved wherever possible so that round-trip writes
// produce minimal diffs.
package envfile
