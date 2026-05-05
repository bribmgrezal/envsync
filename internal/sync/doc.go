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
package sync
