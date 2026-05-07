package envfile

import (
	"testing"
)

func transformEntries() []Entry {
	return []Entry{
		{Key: "db_host", Value: "  localhost  ", LineNo: 1},
		{Key: "db_port", Value: "5432", LineNo: 2},
		{Key: "APP_SECRET", Value: "abc123", LineNo: 3, Masked: true},
		{Key: "STAGING_API_KEY", Value: "key", LineNo: 4},
	}
}

func TestTransform_TrimValues(t *testing.T) {
	opts := DefaultTransformOptions()
	out := Transform(transformEntries(), opts)
	if out[0].Value != "localhost" {
		t.Errorf("expected trimmed value, got %q", out[0].Value)
	}
}

func TestTransform_UppercaseKeys(t *testing.T) {
	opts := DefaultTransformOptions()
	opts.UppercaseKeys = true
	out := Transform(transformEntries(), opts)
	if out[0].Key != "DB_HOST" {
		t.Errorf("expected uppercase key, got %q", out[0].Key)
	}
	if out[1].Key != "DB_PORT" {
		t.Errorf("expected uppercase key, got %q", out[1].Key)
	}
}

func TestTransform_PrefixKeys(t *testing.T) {
	opts := DefaultTransformOptions()
	opts.PrefixKeys = "TEST_"
	out := Transform(transformEntries(), opts)
	if out[0].Key != "TEST_db_host" {
		t.Errorf("expected prefixed key, got %q", out[0].Key)
	}
}

func TestTransform_StripPrefix(t *testing.T) {
	opts := DefaultTransformOptions()
	opts.StripPrefix = "STAGING_"
	out := Transform(transformEntries(), opts)
	if out[3].Key != "API_KEY" {
		t.Errorf("expected stripped key, got %q", out[3].Key)
	}
	// Keys without the prefix should be unchanged.
	if out[0].Key != "db_host" {
		t.Errorf("unexpected key change: %q", out[0].Key)
	}
}

func TestTransform_PreservesMetadata(t *testing.T) {
	opts := DefaultTransformOptions()
	out := Transform(transformEntries(), opts)
	if !out[2].Masked {
		t.Error("expected Masked flag to be preserved")
	}
	if out[2].LineNo != 3 {
		t.Errorf("expected LineNo 3, got %d", out[2].LineNo)
	}
}

func TestTransform_EmptySlice(t *testing.T) {
	out := Transform([]Entry{}, DefaultTransformOptions())
	if len(out) != 0 {
		t.Errorf("expected empty output, got %d entries", len(out))
	}
}

func TestTransform_StripThenPrefix(t *testing.T) {
	opts := DefaultTransformOptions()
	opts.StripPrefix = "STAGING_"
	opts.PrefixKeys = "PROD_"
	out := Transform(transformEntries(), opts)
	if out[3].Key != "PROD_API_KEY" {
		t.Errorf("expected PROD_API_KEY, got %q", out[3].Key)
	}
}
