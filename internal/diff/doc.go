// Package diff implements comparison logic for envsync.
//
// It compares two env maps (typically parsed by the envfile package)
// and classifies each key as one of:
//
//   - match:   key exists in both environments with the same value
//   - missing: key exists in base but is absent from the target
//   - extra:   key exists in target but is not defined in base
//   - changed: key exists in both but values differ
//
// Typical usage:
//
//	base, _ := envfile.Parse("path/to/.env.example")
//	dest, _ := envfile.Parse("path/to/.env.production")
//
//	result := diff.Compare(base, dest)
//	if result.HasDiff() {
//		missing, extra, changed := result.Summary()
//		fmt.Printf("missing=%d extra=%d changed=%d\n", missing, extra, changed)
//	}
//
// The Result.Entries slice contains one Entry per key found in either
// environment, preserving both the base and destination values for
// display or further processing (e.g. secret masking).
package diff
