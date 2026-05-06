package envfile_test

import (
	"testing"

	"github.com/user/envsync/internal/envfile"
)

func TestToMap_BasicConversion(t *testing.T) {
	entries := []envfile.Entry{
		{Key: "FOO", Value: "bar"},
		{Key: "BAZ", Value: "qux"},
	}
	m := envfile.ToMap(entries)
	if m["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %q", m["FOO"])
	}
	if m["BAZ"] != "qux" {
		t.Errorf("expected BAZ=qux, got %q", m["BAZ"])
	}
}

func TestToMap_DuplicateKeepsLast(t *testing.T) {
	entries := []envfile.Entry{
		{Key: "FOO", Value: "first"},
		{Key: "FOO", Value: "second"},
	}
	m := envfile.ToMap(entries)
	if m["FOO"] != "second" {
		t.Errorf("expected last value 'second', got %q", m["FOO"])
	}
}

func TestToMap_EmptySlice(t *testing.T) {
	m := envfile.ToMap(nil)
	if len(m) != 0 {
		t.Errorf("expected empty map, got %d entries", len(m))
	}
}

func TestFromMap_RoundTrip(t *testing.T) {
	orig := map[string]string{"A": "1", "B": "2"}
	entries := envfile.FromMap(orig)
	result := envfile.ToMap(entries)
	for k, v := range orig {
		if result[k] != v {
			t.Errorf("round-trip mismatch for key %q: got %q, want %q", k, result[k], v)
		}
	}
}

func TestFromMap_MaskedAndLineNoUnset(t *testing.T) {
	entries := envfile.FromMap(map[string]string{"X": "y"})
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Masked {
		t.Error("FromMap should not set Masked")
	}
	if entries[0].LineNo != 0 {
		t.Error("FromMap should not set LineNo")
	}
}
