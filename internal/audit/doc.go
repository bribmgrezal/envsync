// Package audit provides functionality for recording and retrieving a history
// of sync operations performed by envsync.
//
// Each audit entry captures the timestamp, source and target file paths,
// the number of changes applied, and the type of operation (e.g. sync, diff).
// Entries are persisted to a JSON file on disk and can be loaded back for
// inspection or reporting purposes.
//
// # Basic Usage
//
// Append a new entry to the audit log:
//
//	 entry := audit.Entry{
//	     Operation: "sync",
//	     Source:    ".env.example",
//	     Target:    ".env",
//	     Changes:   3,
//	 }
//	 if err := audit.Save("/path/to/audit.json", entry); err != nil {
//	     log.Fatal(err)
//	 }
//
// Load all past audit entries:
//
//	 entries, err := audit.Load("/path/to/audit.json")
//	 if err != nil {
//	     log.Fatal(err)
//	 }
//	 for _, e := range entries {
//	     fmt.Println(e.Timestamp, e.Operation, e.Changes)
//	 }
package audit
