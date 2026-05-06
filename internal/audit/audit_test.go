package audit_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/envsync/internal/audit"
)

func TestAppend_AddsEntries(t *testing.T) {
	l := &audit.Log{}
	l.Append("DB_PASSWORD", audit.EventAdded, true)
	l.Append("APP_ENV", audit.EventUpdated, false)

	if len(l.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(l.Entries))
	}
	if l.Entries[0].Key != "DB_PASSWORD" {
		t.Errorf("expected DB_PASSWORD, got %s", l.Entries[0].Key)
	}
	if !l.Entries[0].Masked {
		t.Error("expected first entry to be masked")
	}
	if l.Entries[1].Event != audit.EventUpdated {
		t.Errorf("expected updated event, got %s", l.Entries[1].Event)
	}
}

func TestAppend_SetsTimestamp(t *testing.T) {
	l := &audit.Log{}
	l.Append("KEY", audit.EventAdded, false)

	if l.Entries[0].Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestSaveLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.json")

	original := &audit.Log{}
	original.Append("SECRET_KEY", audit.EventAdded, true)
	original.Append("PORT", audit.EventRemoved, false)

	if err := audit.Save(path, original); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := audit.Load(path)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if len(loaded.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(loaded.Entries))
	}
	if loaded.Entries[0].Key != "SECRET_KEY" {
		t.Errorf("expected SECRET_KEY, got %s", loaded.Entries[0].Key)
	}
	if loaded.Entries[1].Event != audit.EventRemoved {
		t.Errorf("expected removed event, got %s", loaded.Entries[1].Event)
	}
}

func TestLoad_MissingFile_ReturnsError(t *testing.T) {
	_, err := audit.Load("/nonexistent/audit.json")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestSave_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.json")

	l := &audit.Log{}
	l.Append("API_KEY", audit.EventUpdated, true)

	if err := audit.Save(path, l); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error("expected file to exist after Save")
	}
}
