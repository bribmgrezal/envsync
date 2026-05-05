// Package envfile provides functionality for parsing .env files into
// structured representations suitable for diffing and syncing across
// environments.
//
// # Parsing
//
// Use [Parse] to read an .env file from disk. The resulting [EnvFile]
// contains an ordered slice of [Entry] values as well as an index map
// for O(1) key lookups.
//
// Supported syntax:
//   - KEY=VALUE
//   - KEY="quoted value"
//   - KEY=VALUE # inline comment
//   - # full-line comments (skipped)
//   - blank lines (skipped)
//
// Example:
//
//	env, err := envfile.Parse(".env.production")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	for _, e := range env.Entries {
//	    fmt.Printf("%s=%s\n", e.Key, e.Value)
//	}
package envfile
