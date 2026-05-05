package diff_test

import (
	"testing"

	"github.com/user/envsync/internal/diff"
)

func TestCompare_AllMatch(t *testing.T) {
	base := map[string]string{"FOO": "bar", "BAZ": "qux"}
	dest := map[string]string{"FOO": "bar", "BAZ": "qux"}

	r := diff.Compare(base, dest)
	if r.HasDiff() {
		t.Errorf("expected no diff, got entries: %+v", r.Entries)
	}
}

func TestCompare_MissingKey(t *testing.T) {
	base := map[string]string{"FOO": "bar", "MISSING": "val"}
	dest := map[string]string{"FOO": "bar"}

	r := diff.Compare(base, dest)
	missing, extra, changed := r.Summary()

	if missing != 1 {
		t.Errorf("expected 1 missing, got %d", missing)
	}
	if extra != 0 || changed != 0 {
		t.Errorf("unexpected extra=%d changed=%d", extra, changed)
	}
}

func TestCompare_ExtraKey(t *testing.T) {
	base := map[string]string{"FOO": "bar"}
	dest := map[string]string{"FOO": "bar", "EXTRA": "val"}

	r := diff.Compare(base, dest)
	_, extra, _ := r.Summary()

	if extra != 1 {
		t.Errorf("expected 1 extra, got %d", extra)
	}
}

func TestCompare_ChangedValue(t *testing.T) {
	base := map[string]string{"FOO": "original"}
	dest := map[string]string{"FOO": "modified"}

	r := diff.Compare(base, dest)
	_, _, changed := r.Summary()

	if changed != 1 {
		t.Errorf("expected 1 changed, got %d", changed)
	}

	for _, e := range r.Entries {
		if e.Key == "FOO" {
			if e.BaseValue != "original" {
				t.Errorf("expected BaseValue=original, got %q", e.BaseValue)
			}
			if e.DestValue != "modified" {
				t.Errorf("expected DestValue=modified, got %q", e.DestValue)
			}
		}
	}
}

func TestCompare_EmptyBoth(t *testing.T) {
	r := diff.Compare(map[string]string{}, map[string]string{})
	if r.HasDiff() {
		t.Error("expected no diff for two empty maps")
	}
	if len(r.Entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(r.Entries))
	}
}

func TestCompare_MixedDiffs(t *testing.T) {
	base := map[string]string{"A": "1", "B": "2", "C": "3"}
	dest := map[string]string{"A": "1", "B": "changed", "D": "4"}

	r := diff.Compare(base, dest)
	missing, extra, changed := r.Summary()

	if missing != 1 {
		t.Errorf("expected 1 missing (C), got %d", missing)
	}
	if extra != 1 {
		t.Errorf("expected 1 extra (D), got %d", extra)
	}
	if changed != 1 {
		t.Errorf("expected 1 changed (B), got %d", changed)
	}
}
