package envfile

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func patchBase() []Entry {
	return []Entry{
		{Key: "APP_ENV", Value: "development", LineNo: 1},
		{Key: "DB_HOST", Value: "localhost", LineNo: 2},
		{Key: "DB_PORT", Value: "5432", LineNo: 3},
	}
}

func TestPatch_OverwritesExistingKey(t *testing.T) {
	patch := []Entry{{Key: "DB_HOST", Value: "prod-db.internal"}}
	got, err := Patch(patchBase(), patch, DefaultPatchOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got[1].Value != "prod-db.internal" {
		t.Errorf("expected updated value, got %q", got[1].Value)
	}
	// LineNo must be preserved from base.
	if got[1].LineNo != 2 {
		t.Errorf("expected LineNo 2, got %d", got[1].LineNo)
	}
}

func TestPatch_AppendsNewKey_WhenAllowNew(t *testing.T) {
	opts := DefaultPatchOptions()
	opts.AllowNew = true
	patch := []Entry{{Key: "REDIS_URL", Value: "redis://localhost:6379"}}
	got, err := Patch(patchBase(), patch, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 4 {
		t.Fatalf("expected 4 entries, got %d", len(got))
	}
	if got[3].Key != "REDIS_URL" {
		t.Errorf("expected REDIS_URL appended, got %q", got[3].Key)
	}
}

func TestPatch_SkipsNewKey_WhenAllowNewFalse(t *testing.T) {
	opts := DefaultPatchOptions()
	opts.AllowNew = false
	patch := []Entry{{Key: "NEW_KEY", Value: "value"}}
	got, err := Patch(patchBase(), patch, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 3 {
		t.Errorf("expected 3 entries, got %d", len(got))
	}
}

func TestPatch_FailOnMissing_ReturnsError(t *testing.T) {
	opts := DefaultPatchOptions()
	opts.AllowNew = false
	opts.FailOnMissing = true
	patch := []Entry{{Key: "GHOST_KEY", Value: "x"}}
	_, err := Patch(patchBase(), patch, opts)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestPatch_PreservesOrder(t *testing.T) {
	patch := []Entry{
		{Key: "DB_PORT", Value: "5433"},
		{Key: "APP_ENV", Value: "production"},
	}
	got, err := Patch(patchBase(), patch, DefaultPatchOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []Entry{
		{Key: "APP_ENV", Value: "production", LineNo: 1},
		{Key: "DB_HOST", Value: "localhost", LineNo: 2},
		{Key: "DB_PORT", Value: "5433", LineNo: 3},
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("order mismatch (-want +got):\n%s", diff)
	}
}

func TestPatch_DoesNotMutateBase(t *testing.T) {
	base := patchBase()
	origVal := base[0].Value
	patch := []Entry{{Key: "APP_ENV", Value: "staging"}}
	_, _ = Patch(base, patch, DefaultPatchOptions())
	if base[0].Value != origVal {
		t.Errorf("base was mutated: got %q, want %q", base[0].Value, origVal)
	}
}
