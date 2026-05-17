// Package envfile provides utilities for parsing, transforming, and
// manipulating .env files.
//
// # Coalesce
//
// Coalesce merges multiple ordered slices of Entry into a single deduplicated
// slice. It is useful when you want to combine values from several environment
// sources — for example a base .env, a per-environment override, and a
// secrets file — while controlling which source takes precedence.
//
// Resolution rules
//
//   - By default (PreferFirst: false) the LAST source that provides a
//     non-empty value for a given key wins. This mirrors the typical shell
//     behaviour where later exports override earlier ones.
//
//   - When PreferFirst is true the FIRST source that provides a non-empty
//     value wins. Subsequent sources are ignored for that key. This is
//     useful when a local developer override should never be clobbered by
//     a shared defaults file.
//
//   - SkipEmpty (default true) treats empty-string values as absent so that
//     a source with a placeholder empty value does not shadow a real value
//     in another source.
//
// The relative insertion order of keys is preserved: a key retains the
// position it first appeared at, even if its final value comes from a later
// source.
package envfile
