package envfile

import (
	"testing"
)

func baseEntries() []Entry {
	return []Entry{
		{Key: "APP_ENV", Value: "production"},
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "DB_PASS", Value: "secret"},
	}
}

func TestCompareEntries_AllSame(t *testing.T) {
	base := baseEntries()
	other := baseEntries()
	results := CompareEntries(base, other)
	for _, d := range results {
		if d.Status != DiffSame {
			t.Errorf("expected all same, got %s for key %s", d.Status, d.Key)
		}
	}
}

func TestCompareEntries_Added(t *testing.T) {
	base := baseEntries()
	other := append(baseEntries(), Entry{Key: "NEW_KEY", Value: "new"})
	results := CompareEntries(base, other)
	found := false
	for _, d := range results {
		if d.Key == "NEW_KEY" && d.Status == DiffAdded {
			found = true
		}
	}
	if !found {
		t.Error("expected NEW_KEY to be marked as added")
	}
}

func TestCompareEntries_Removed(t *testing.T) {
	base := baseEntries()
	other := []Entry{
		{Key: "APP_ENV", Value: "production"},
		{Key: "DB_HOST", Value: "localhost"},
	}
	results := CompareEntries(base, other)
	for _, d := range results {
		if d.Key == "DB_PASS" && d.Status != DiffRemoved {
			t.Errorf("expected DB_PASS removed, got %s", d.Status)
		}
	}
}

func TestCompareEntries_Changed(t *testing.T) {
	base := baseEntries()
	other := []Entry{
		{Key: "APP_ENV", Value: "staging"},
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "DB_PASS", Value: "secret"},
	}
	results := CompareEntries(base, other)
	for _, d := range results {
		if d.Key == "APP_ENV" {
			if d.Status != DiffChanged {
				t.Errorf("expected APP_ENV changed, got %s", d.Status)
			}
			if d.BaseVal != "production" || d.OtherVal != "staging" {
				t.Errorf("unexpected values: base=%s other=%s", d.BaseVal, d.OtherVal)
			}
		}
	}
}

func TestCompareEntries_SortedByKey(t *testing.T) {
	base := []Entry{{Key: "Z_KEY", Value: "z"}, {Key: "A_KEY", Value: "a"}}
	other := []Entry{{Key: "M_KEY", Value: "m"}}
	results := CompareEntries(base, other)
	for i := 1; i < len(results); i++ {
		if results[i-1].Key > results[i].Key {
			t.Errorf("results not sorted: %s > %s", results[i-1].Key, results[i].Key)
		}
	}
}

func TestDiffStatus_String(t *testing.T) {
	cases := map[DiffStatus]string{
		DiffAdded:   "added",
		DiffRemoved: "removed",
		DiffChanged: "changed",
		DiffSame:    "same",
	}
	for status, want := range cases {
		if got := status.String(); got != want {
			t.Errorf("DiffStatus(%d).String() = %q, want %q", status, got, want)
		}
	}
}
