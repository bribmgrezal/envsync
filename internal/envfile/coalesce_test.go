package envfile

import (
	"testing"
)

func coalesceEntry(key, value string) Entry {
	return Entry{Key: key, Value: value}
}

func TestCoalesce_LastSourceWins(t *testing.T) {
	a := []Entry{coalesceEntry("HOST", "localhost"), coalesceEntry("PORT", "3000")}
	b := []Entry{coalesceEntry("HOST", "prod.example.com")}

	result := Coalesce([][]Entry{a, b}, DefaultCoalesceOptions())

	m := ToMap(result)
	if m["HOST"] != "prod.example.com" {
		t.Errorf("expected prod.example.com, got %s", m["HOST"])
	}
	if m["PORT"] != "3000" {
		t.Errorf("expected 3000, got %s", m["PORT"])
	}
}

func TestCoalesce_PreferFirst(t *testing.T) {
	a := []Entry{coalesceEntry("DB", "postgres://local")}
	b := []Entry{coalesceEntry("DB", "postgres://remote")}

	opts := DefaultCoalesceOptions()
	opts.PreferFirst = true

	result := Coalesce([][]Entry{a, b}, opts)

	m := ToMap(result)
	if m["DB"] != "postgres://local" {
		t.Errorf("expected postgres://local, got %s", m["DB"])
	}
}

func TestCoalesce_SkipEmpty_FillsFromLaterSource(t *testing.T) {
	a := []Entry{coalesceEntry("SECRET", "")}
	b := []Entry{coalesceEntry("SECRET", "abc123")}

	opts := DefaultCoalesceOptions() // SkipEmpty true by default

	result := Coalesce([][]Entry{a, b}, opts)

	m := ToMap(result)
	if m["SECRET"] != "abc123" {
		t.Errorf("expected abc123, got %s", m["SECRET"])
	}
}

func TestCoalesce_SkipEmptyFalse_EmptyValueWins(t *testing.T) {
	a := []Entry{coalesceEntry("KEY", "original")}
	b := []Entry{coalesceEntry("KEY", "")}

	opts := DefaultCoalesceOptions()
	opts.SkipEmpty = false

	result := Coalesce([][]Entry{a, b}, opts)

	m := ToMap(result)
	if m["KEY"] != "" {
		t.Errorf("expected empty string, got %s", m["KEY"])
	}
}

func TestCoalesce_PreservesInsertionOrder(t *testing.T) {
	a := []Entry{coalesceEntry("A", "1"), coalesceEntry("B", "2")}
	b := []Entry{coalesceEntry("C", "3"), coalesceEntry("A", "99")}

	result := Coalesce([][]Entry{a, b}, DefaultCoalesceOptions())

	if len(result) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(result))
	}
	// A should appear at index 0 (its original position)
	if result[0].Key != "A" {
		t.Errorf("expected A at position 0, got %s", result[0].Key)
	}
	if result[0].Value != "99" {
		t.Errorf("expected value 99, got %s", result[0].Value)
	}
}

func TestCoalesce_EmptySources_ReturnsEmpty(t *testing.T) {
	result := Coalesce([][]Entry{}, DefaultCoalesceOptions())
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d entries", len(result))
	}
}

func TestCoalesce_SingleSource_PassThrough(t *testing.T) {
	src := []Entry{coalesceEntry("X", "1"), coalesceEntry("Y", "2")}
	result := Coalesce([][]Entry{src}, DefaultCoalesceOptions())
	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
}
