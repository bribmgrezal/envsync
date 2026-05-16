package envfile

import (
	"testing"
)

func promoteEntries() []Entry {
	return []Entry{
		{Key: "DEV_DB_HOST", Value: "localhost", LineNo: 1},
		{Key: "DEV_DB_PORT", Value: "5432", LineNo: 2},
		{Key: "APP_NAME", Value: "envsync", LineNo: 3},
	}
}

func TestPromote_StripsFromPrefixAndAppliesTo(t *testing.T) {
	opts := DefaultPromoteOptions()
	opts.FromPrefix = "DEV_"
	opts.ToPrefix = "PROD_"

	out, err := Promote(promoteEntries(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// APP_NAME kept + 2 promoted entries (DEV_ entries replaced)
	if len(out) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(out))
	}
	if out[0].Key != "APP_NAME" {
		t.Errorf("expected APP_NAME, got %s", out[0].Key)
	}
	if out[1].Key != "PROD_DB_HOST" {
		t.Errorf("expected PROD_DB_HOST, got %s", out[1].Key)
	}
	if out[2].Key != "PROD_DB_PORT" {
		t.Errorf("expected PROD_DB_PORT, got %s", out[2].Key)
	}
}

func TestPromote_KeepOriginal_RetainsBothEntries(t *testing.T) {
	opts := DefaultPromoteOptions()
	opts.FromPrefix = "DEV_"
	opts.ToPrefix = "PROD_"
	opts.KeepOriginal = true

	out, err := Promote(promoteEntries(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// APP_NAME + 2*(original+promoted) = 5
	if len(out) != 5 {
		t.Fatalf("expected 5 entries, got %d", len(out))
	}
}

func TestPromote_NoFromPrefix_PrependsToPrefixAll(t *testing.T) {
	opts := DefaultPromoteOptions()
	opts.ToPrefix = "EXPORT_"

	out, err := Promote(promoteEntries(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, e := range out {
		if len(e.Key) < len("EXPORT_") || e.Key[:7] != "EXPORT_" {
			t.Errorf("expected EXPORT_ prefix, got %s", e.Key)
		}
	}
}

func TestPromote_FailOnMissing_ReturnsError(t *testing.T) {
	opts := DefaultPromoteOptions()
	opts.FromPrefix = "STAGING_"
	opts.FailOnMissing = true

	_, err := Promote(promoteEntries(), opts)
	if err == nil {
		t.Fatal("expected error for missing prefix, got nil")
	}
}

func TestPromote_EmptyEntries_ReturnsEmpty(t *testing.T) {
	opts := DefaultPromoteOptions()
	opts.FromPrefix = "DEV_"
	opts.ToPrefix = "PROD_"

	out, err := Promote([]Entry{}, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 0 {
		t.Errorf("expected empty slice, got %d entries", len(out))
	}
}

func TestPromote_PromotedEntry_HasZeroLineNo(t *testing.T) {
	opts := DefaultPromoteOptions()
	opts.FromPrefix = "DEV_"
	opts.ToPrefix = "PROD_"

	out, _ := Promote(promoteEntries(), opts)
	for _, e := range out {
		if e.Key == "PROD_DB_HOST" && e.LineNo != 0 {
			t.Errorf("promoted entry should have LineNo=0, got %d", e.LineNo)
		}
	}
}
