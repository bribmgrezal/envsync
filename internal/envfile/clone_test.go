package envfile

import (
	"testing"
)

var cloneEntries = []Entry{
	{Key: "DB_HOST", Value: "localhost", LineNo: 1},
	{Key: "DB_PORT", Value: "5432", LineNo: 2},
	{Key: "API_KEY", Value: "secret", LineNo: 3},
}

func TestClone_DuplicatesKeyUnderNewName(t *testing.T) {
	opts := DefaultCloneOptions()
	opts.KeyMap = map[string][]string{"DB_HOST": {"DATABASE_HOST"}}

	out, err := Clone(cloneEntries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 4 {
		t.Fatalf("expected 4 entries, got %d", len(out))
	}
	if out[1].Key != "DATABASE_HOST" || out[1].Value != "localhost" {
		t.Errorf("expected cloned entry DATABASE_HOST=localhost, got %s=%s", out[1].Key, out[1].Value)
	}
	if out[1].LineNo != 0 {
		t.Errorf("expected LineNo reset to 0, got %d", out[1].LineNo)
	}
}

func TestClone_MultipleDestinations(t *testing.T) {
	opts := DefaultCloneOptions()
	opts.KeyMap = map[string][]string{"DB_PORT": {"PG_PORT", "POSTGRES_PORT"}}

	out, err := Clone(cloneEntries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// 3 originals + 2 clones = 5
	if len(out) != 5 {
		t.Fatalf("expected 5 entries, got %d", len(out))
	}
}

func TestClone_KeepOriginalFalse_DropsSrc(t *testing.T) {
	opts := DefaultCloneOptions()
	opts.KeepOriginal = false
	opts.KeyMap = map[string][]string{"DB_HOST": {"DATABASE_HOST"}}

	out, err := Clone(cloneEntries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, e := range out {
		if e.Key == "DB_HOST" {
			t.Errorf("expected DB_HOST to be removed, but found it")
		}
	}
}

func TestClone_MissingKey_ReturnsError(t *testing.T) {
	opts := DefaultCloneOptions()
	opts.KeyMap = map[string][]string{"MISSING_KEY": {"OTHER_KEY"}}

	_, err := Clone(cloneEntries, opts)
	if err == nil {
		t.Fatal("expected error for missing source key, got nil")
	}
}

func TestClone_SkipMissing_NoError(t *testing.T) {
	opts := DefaultCloneOptions()
	opts.SkipMissing = true
	opts.KeyMap = map[string][]string{"MISSING_KEY": {"OTHER_KEY"}}

	out, err := Clone(cloneEntries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != len(cloneEntries) {
		t.Errorf("expected %d entries unchanged, got %d", len(cloneEntries), len(out))
	}
}

func TestClone_EmptyKeyMap_ReturnsOriginal(t *testing.T) {
	opts := DefaultCloneOptions()

	out, err := Clone(cloneEntries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != len(cloneEntries) {
		t.Errorf("expected %d entries, got %d", len(cloneEntries), len(out))
	}
}
