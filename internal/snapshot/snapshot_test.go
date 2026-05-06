package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourusername/envsync/internal/envfile"
	"github.com/yourusername/envsync/internal/snapshot"
)

func sampleEntries() []envfile.Entry {
	return []envfile.Entry{
		{Key: "APP_ENV", Value: "production", LineNo: 1},
		{Key: "DB_PASSWORD", Value: "secret123", LineNo: 2},
		{Key: "PORT", Value: "8080", LineNo: 3},
	}
}

func TestSave_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	err := snapshot.Save(path, ".env", sampleEntries())
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatal("expected snapshot file to exist")
	}
}

func TestLoad_ReturnsEntries(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")
	entries := sampleEntries()

	if err := snapshot.Save(path, ".env.production", entries); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	snap, err := snapshot.Load(path)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if snap.Source != ".env.production" {
		t.Errorf("Source = %q, want %q", snap.Source, ".env.production")
	}

	if len(snap.Entries) != len(entries) {
		t.Fatalf("len(Entries) = %d, want %d", len(snap.Entries), len(entries))
	}

	for i, e := range entries {
		if snap.Entries[i].Key != e.Key || snap.Entries[i].Value != e.Value {
			t.Errorf("entry[%d] = %+v, want %+v", i, snap.Entries[i], e)
		}
	}
}

func TestLoad_SetsCreatedAt(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")
	before := time.Now().UTC().Add(-time.Second)

	if err := snapshot.Save(path, ".env", sampleEntries()); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	snap, err := snapshot.Load(path)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if snap.CreatedAt.Before(before) {
		t.Errorf("CreatedAt = %v, expected after %v", snap.CreatedAt, before)
	}
}

func TestLoad_MissingFile_ReturnsError(t *testing.T) {
	_, err := snapshot.Load("/nonexistent/path/snap.json")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoad_InvalidJSON_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")

	if err := os.WriteFile(path, []byte("not-json{"), 0600); err != nil {
		t.Fatalf("WriteFile error = %v", err)
	}

	_, err := snapshot.Load(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}
