// Package snapshot provides functionality for saving and loading point-in-time
// captures of environment file entries.
//
// A snapshot records the full set of key-value entries from a parsed .env file
// along with metadata such as the source filename and the time the snapshot was
// created. Snapshots are persisted as JSON files and can be reloaded for
// comparison, auditing, or rollback purposes.
//
// Basic usage:
//
//	entries, _ := envfile.Parse(".env.production")
//	snapshot.Save("snapshots/prod.json", ".env.production", entries)
//
//	snap, _ := snapshot.Load("snapshots/prod.json")
//	fmt.Println(snap.Source, snap.CreatedAt)
package snapshot
