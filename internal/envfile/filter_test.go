package envfile

import (
	"testing"
)

func sampleFilterEntries() []Entry {
	return []Entry{
		{Key: "APP_NAME", Value: "envsync"},
		{Key: "APP_ENV", Value: "production"},
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "DB_PASSWORD", Value: ""},
		{Key: "SECRET_KEY", Value: "abc123"},
	}
}

func TestFilter_NoOptions_ReturnsAll(t *testing.T) {
	entries := sampleFilterEntries()
	got := Filter(entries, DefaultFilterOptions())
	if len(got) != len(entries) {
		t.Fatalf("expected %d entries, got %d", len(entries), len(got))
	}
}

func TestFilter_Prefix_KeepsMatching(t *testing.T) {
	got := Filter(sampleFilterEntries(), FilterOptions{Prefix: "APP_"})
	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(got))
	}
	for _, e := range got {
		if e.Key != "APP_NAME" && e.Key != "APP_ENV" {
			t.Errorf("unexpected key %q in result", e.Key)
		}
	}
}

func TestFilter_ExcludeKeys_DropsListed(t *testing.T) {
	opts := FilterOptions{ExcludeKeys: []string{"DB_PASSWORD", "SECRET_KEY"}}
	got := Filter(sampleFilterEntries(), opts)
	for _, e := range got {
		if e.Key == "DB_PASSWORD" || e.Key == "SECRET_KEY" {
			t.Errorf("excluded key %q still present", e.Key)
		}
	}
	if len(got) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(got))
	}
}

func TestFilter_SkipEmpty_OmitsBlankValues(t *testing.T) {
	got := Filter(sampleFilterEntries(), FilterOptions{SkipEmpty: true})
	for _, e := range got {
		if e.Value == "" {
			t.Errorf("entry with empty value %q not removed", e.Key)
		}
	}
	if len(got) != 4 {
		t.Fatalf("expected 4 entries after skipping empty, got %d", len(got))
	}
}

func TestFilter_CombinedOptions(t *testing.T) {
	opts := FilterOptions{
		Prefix:      "DB_",
		SkipEmpty:   true,
		ExcludeKeys: []string{"DB_HOST"},
	}
	got := Filter(sampleFilterEntries(), opts)
	// DB_HOST excluded, DB_PASSWORD empty -> 0 results
	if len(got) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(got))
	}
}

func TestFilter_EmptyInput_ReturnsEmpty(t *testing.T) {
	got := Filter([]Entry{}, FilterOptions{Prefix: "APP_"})
	if len(got) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(got))
	}
}

func TestFilter_Prefix_NoMatch_ReturnsEmpty(t *testing.T) {
	got := Filter(sampleFilterEntries(), FilterOptions{Prefix: "UNKNOWN_"})
	if len(got) != 0 {
		t.Fatalf("expected 0 entries for non-matching prefix, got %d", len(got))
	}
}
