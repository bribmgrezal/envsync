package envfile

import (
	"testing"
)

func renameEntries() []Entry {
	return []Entry{
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "DB_PORT", Value: "5432"},
		{Key: "APP_SECRET", Value: "abc123"},
	}
}

func TestRename_RenamesMatchingKey(t *testing.T) {
	opts := DefaultRenameOptions()
	opts.Mapping = map[string]string{"DB_HOST": "DATABASE_HOST"}

	out, err := Rename(renameEntries(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0].Key != "DATABASE_HOST" {
		t.Errorf("expected DATABASE_HOST, got %q", out[0].Key)
	}
}

func TestRename_LeavesUnmappedKeysAlone(t *testing.T) {
	opts := DefaultRenameOptions()
	opts.Mapping = map[string]string{"DB_HOST": "DATABASE_HOST"}

	out, err := Rename(renameEntries(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[1].Key != "DB_PORT" {
		t.Errorf("expected DB_PORT unchanged, got %q", out[1].Key)
	}
	if out[2].Key != "APP_SECRET" {
		t.Errorf("expected APP_SECRET unchanged, got %q", out[2].Key)
	}
}

func TestRename_PreservesValues(t *testing.T) {
	opts := DefaultRenameOptions()
	opts.Mapping = map[string]string{"DB_PORT": "DATABASE_PORT"}

	out, err := Rename(renameEntries(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[1].Value != "5432" {
		t.Errorf("expected value 5432, got %q", out[1].Value)
	}
}

func TestRename_SkipMissing_True_NoError(t *testing.T) {
	opts := DefaultRenameOptions()
	opts.SkipMissing = true
	opts.Mapping = map[string]string{"NONEXISTENT_KEY": "NEW_KEY"}

	_, err := Rename(renameEntries(), opts)
	if err != nil {
		t.Errorf("expected no error with SkipMissing=true, got: %v", err)
	}
}

func TestRename_SkipMissing_False_ReturnsError(t *testing.T) {
	opts := DefaultRenameOptions()
	opts.SkipMissing = false
	opts.Mapping = map[string]string{"NONEXISTENT_KEY": "NEW_KEY"}

	_, err := Rename(renameEntries(), opts)
	if err == nil {
		t.Error("expected error for missing key with SkipMissing=false, got nil")
	}
}

func TestRename_EmptyMapping_ReturnsOriginal(t *testing.T) {
	opts := DefaultRenameOptions()

	out, err := Rename(renameEntries(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 3 {
		t.Errorf("expected 3 entries, got %d", len(out))
	}
	if out[0].Key != "DB_HOST" {
		t.Errorf("expected DB_HOST, got %q", out[0].Key)
	}
}

func TestRename_MultipleKeys(t *testing.T) {
	opts := DefaultRenameOptions()
	opts.Mapping = map[string]string{
		"DB_HOST":    "DATABASE_HOST",
		"APP_SECRET": "APPLICATION_SECRET",
	}

	out, err := Rename(renameEntries(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0].Key != "DATABASE_HOST" {
		t.Errorf("expected DATABASE_HOST, got %q", out[0].Key)
	}
	if out[2].Key != "APPLICATION_SECRET" {
		t.Errorf("expected APPLICATION_SECRET, got %q", out[2].Key)
	}
}
