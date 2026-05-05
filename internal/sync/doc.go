// Package sync implements the synchronisation step of envsync.
//
// After a diff has been computed between a source and target environment
// file (see package diff), Apply can be used to merge the desired changes
// back into the target map according to a set of configurable Options:
//
//	// Build a merged map from diff results.
//	updated, summary := sync.Apply(results, targetEnv, sync.DefaultOptions())
//	fmt.Println(summary)
//
// The returned map can then be serialised back to disk by the caller.
// sync intentionally does not perform any I/O so that it remains easy
// to test and compose with other packages.
//
// # Options
//
// DefaultOptions returns a sensible baseline configuration. Individual
// fields can be overridden before passing the Options value to Apply:
//
//	opts := sync.DefaultOptions()
//	opts.OverwriteExisting = true
//	updated, summary := sync.Apply(results, targetEnv, opts)
//
// # Summary
//
// Apply also returns a human-readable Summary that describes how many
// keys were added, updated, or skipped during the merge. This is
// suitable for printing directly to the user or logging for auditing.
package sync
