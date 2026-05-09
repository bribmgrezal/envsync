package envfile

import (
	"testing"
)

var sortEntries = []Entry{
	{Key: "ZEBRA", Value: "z"},
	{Key: "APP_NAME", Value: "envsync"},
	{Key: "DB_HOST", Value: "localhost"},
	{Key: "APP_ENV", Value: "prod"},
	{Key: "DB_PORT", Value: "5432"},
	{Key: "ALPHA", Value: "a"},
}

func TestSort_Alpha(t *testing.T) {
	result := Sort(sortEntries, SortOptions{Order: SortAlpha})
	expected := []string{"ALPHA", "APP_ENV", "APP_NAME", "DB_HOST", "DB_PORT", "ZEBRA"}
	assertKeyOrder(t, result, expected)
}

func TestSort_AlphaDesc(t *testing.T) {
	result := Sort(sortEntries, SortOptions{Order: SortAlphaDesc})
	expected := []string{"ZEBRA", "DB_PORT", "DB_HOST", "APP_NAME", "APP_ENV", "ALPHA"}
	assertKeyOrder(t, result, expected)
}

func TestSort_ByGroup(t *testing.T) {
	result := Sort(sortEntries, SortOptions{Order: SortByGroup})
	// Groups: ALPHA (no _), APP_, DB_, ZEBRA (no _)
	// Within group: alphabetical
	expected := []string{"ALPHA", "APP_ENV", "APP_NAME", "DB_HOST", "DB_PORT", "ZEBRA"}
	assertKeyOrder(t, result, expected)
}

func TestSort_DoesNotMutateOriginal(t *testing.T) {
	orig := []Entry{
		{Key: "Z", Value: "z"},
		{Key: "A", Value: "a"},
	}
	_ = Sort(orig, DefaultSortOptions())
	if orig[0].Key != "Z" {
		t.Errorf("original slice was mutated: got %s, want Z", orig[0].Key)
	}
}

func TestSort_EmptySlice(t *testing.T) {
	result := Sort([]Entry{}, DefaultSortOptions())
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d entries", len(result))
	}
}

func TestDefaultSortOptions(t *testing.T) {
	opts := DefaultSortOptions()
	if opts.Order != SortAlpha {
		t.Errorf("expected SortAlpha, got %v", opts.Order)
	}
}

func assertKeyOrder(t *testing.T, entries []Entry, expected []string) {
	t.Helper()
	if len(entries) != len(expected) {
		t.Fatalf("length mismatch: got %d, want %d", len(entries), len(expected))
	}
	for i, e := range entries {
		if e.Key != expected[i] {
			t.Errorf("index %d: got %q, want %q", i, e.Key, expected[i])
		}
	}
}
