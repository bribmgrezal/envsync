// Package audit provides functionality to record and retrieve a history
// of sync operations performed by envsync, enabling traceability of
// environment changes over time.
package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// EventType describes the kind of sync action that was recorded.
type EventType string

const (
	EventAdded   EventType = "added"
	EventUpdated EventType = "updated"
	EventRemoved EventType = "removed"
)

// Entry represents a single audited change to an environment key.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Key       string    `json:"key"`
	Event     EventType `json:"event"`
	Masked    bool      `json:"masked"`
}

// Log holds an ordered list of audit entries.
type Log struct {
	Entries []Entry `json:"entries"`
}

// Append adds a new audit entry to the log.
func (l *Log) Append(key string, event EventType, masked bool) {
	l.Entries = append(l.Entries, Entry{
		Timestamp: time.Now().UTC(),
		Key:       key,
		Event:     event,
		Masked:    masked,
	})
}

// Save writes the audit log as JSON to the given file path.
// If the file already exists it is overwritten.
func Save(path string, l *Log) error {
	data, err := json.MarshalIndent(l, "", "  ")
	if err != nil {
		return fmt.Errorf("audit: marshal: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("audit: write %s: %w", path, err)
	}
	return nil
}

// Load reads a previously saved audit log from the given file path.
func Load(path string) (*Log, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("audit: read %s: %w", path, err)
	}
	var l Log
	if err := json.Unmarshal(data, &l); err != nil {
		return nil, fmt.Errorf("audit: unmarshal: %w", err)
	}
	return &l, nil
}
