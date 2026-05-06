package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/yourusername/envsync/internal/envfile"
)

// Snapshot represents a point-in-time capture of an env file's entries.
type Snapshot struct {
	CreatedAt time.Time          `json:"created_at"`
	Source    string             `json:"source"`
	Entries   []envfile.Entry    `json:"entries"`
}

// Save writes a snapshot of the given entries to the specified file path as JSON.
func Save(path, source string, entries []envfile.Entry) error {
	snap := Snapshot{
		CreatedAt: time.Now().UTC(),
		Source:    source,
		Entries:   entries,
	}

	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshot: marshal failed: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("snapshot: write failed: %w", err)
	}

	return nil
}

// Load reads a snapshot from the given file path.
func Load(path string) (*Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("snapshot: read failed: %w", err)
	}

	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return nil, fmt.Errorf("snapshot: parse failed: %w", err)
	}

	return &snap, nil
}
