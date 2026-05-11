package envfile

import (
	"testing"
)

func dupeEntries() []Entry {
	return []Entry{
		{Key: "HOST", Value: "localhost", LineNo: 1},
		{Key: "PORT", Value: "3000", LineNo: 2},
		{Key: "HOST", Value: "remotehost", LineNo: 3},
		{Key: "DEBUG", Value: "true", LineNo: 4},
		{Key: "PORT", Value: "4000", LineNo: 5},
	}
}

func TestDedupe_KeepsLastByDefault(t *testing.T) {
	result := Dedupe(dupeEntries(), DefaultDedupeOptions())
	if len(result.Entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(result.Entries))
	}
	assertEntry(t, result.Entries, "HOST", "remotehost")
	assertEntry(t, result.Entries, "PORT", "4000")
	assertEntry(t, result.Entries, "DEBUG", "true")
}

func TestDedupe_KeepFirst(t *testing.T) {
	opts := DefaultDedupeOptions()
	opts.KeepFirst = true
	result := Dedupe(dupeEntries(), opts)
	if len(result.Entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(result.Entries))
	}
	assertEntry(t, result.Entries, "HOST", "localhost")
	assertEntry(t, result.Entries, "PORT", "3000")
}

func TestDedupe_ReportsDuplicates(t *testing.T) {
	opts := DefaultDedupeOptions()
	opts.ReportDuplicates = true
	result := Dedupe(dupeEntries(), opts)
	if len(result.Duplicates) != 2 {
		t.Fatalf("expected 2 duplicate keys reported, got %d", len(result.Duplicates))
	}
	dupeMap := make(map[string]bool)
	for _, k := range result.Duplicates {
		dupeMap[k] = true
	}
	if !dupeMap["HOST"] || !dupeMap["PORT"] {
		t.Errorf("expected HOST and PORT in duplicates, got %v", result.Duplicates)
	}
}

func TestDedupe_NoDuplicates_ReturnsAll(t *testing.T) {
	entries := []Entry{
		{Key: "A", Value: "1", LineNo: 1},
		{Key: "B", Value: "2", LineNo: 2},
	}
	result := Dedupe(entries, DefaultDedupeOptions())
	if len(result.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result.Entries))
	}
}

func TestDedupe_EmptyInput(t *testing.T) {
	result := Dedupe([]Entry{}, DefaultDedupeOptions())
	if len(result.Entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(result.Entries))
	}
}

func TestDedupe_PreservesOrder(t *testing.T) {
	opts := DefaultDedupeOptions()
	result := Dedupe(dupeEntries(), opts)
	// HOST appears at line 3 (last), PORT at line 5 (last), DEBUG at line 4.
	// Expected order by first-seen position of surviving entry: HOST(1→3), PORT(2→5), DEBUG(4).
	keys := []string{"HOST", "PORT", "DEBUG"}
	for i, e := range result.Entries {
		if e.Key != keys[i] {
			t.Errorf("position %d: expected key %q, got %q", i, keys[i], e.Key)
		}
	}
}

// assertEntry is a helper to find a key in entries and check its value.
func assertEntry(t *testing.T, entries []Entry, key, wantValue string) {
	t.Helper()
	for _, e := range entries {
		if e.Key == key {
			if e.Value != wantValue {
				t.Errorf("key %q: expected value %q, got %q", key, wantValue, e.Value)
			}
			return
		}
	}
	t.Errorf("key %q not found in entries", key)
}
