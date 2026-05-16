package envfile

import (
	"testing"
)

var pivotEntries = []Entry{
	{Key: "APP_NAME", Value: "envsync"},
	{Key: "APP_VERSION", Value: "1.0"},
	{Key: "KEY", Value: "DB_HOST"},
	{Key: "VALUE", Value: "localhost"},
}

func TestPivot_CollapsesPairIntoEntry(t *testing.T) {
	opts := DefaultPivotOptions()
	out, err := Pivot(pivotEntries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Should contain APP_NAME, APP_VERSION, and the pivoted DB_HOST.
	if len(out) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(out))
	}
	last := out[len(out)-1]
	if last.Key != "DB_HOST" {
		t.Errorf("expected pivoted key DB_HOST, got %q", last.Key)
	}
	if last.Value != "localhost" {
		t.Errorf("expected pivoted value localhost, got %q", last.Value)
	}
}

func TestPivot_PassThroughWhenNoPivotColumns(t *testing.T) {
	entries := []Entry{
		{Key: "FOO", Value: "bar"},
		{Key: "BAZ", Value: "qux"},
	}
	out, err := Pivot(entries, DefaultPivotOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 entries unchanged, got %d", len(out))
	}
}

func TestPivot_UppercaseKey(t *testing.T) {
	entries := []Entry{
		{Key: "key", Value: "api_token"},
		{Key: "value", Value: "secret123"},
	}
	opts := DefaultPivotOptions()
	opts.KeyColumn = "key"
	opts.ValColumn = "value"
	out, err := Pivot(entries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("expected 1 pivoted entry, got %d", len(out))
	}
	if out[0].Key != "API_TOKEN" {
		t.Errorf("expected uppercase key API_TOKEN, got %q", out[0].Key)
	}
}

func TestPivot_FailOnMissing_MissingValColumn(t *testing.T) {
	entries := []Entry{
		{Key: "KEY", Value: "DB_HOST"},
	}
	opts := DefaultPivotOptions()
	opts.FailOnMissing = true
	_, err := Pivot(entries, opts)
	if err == nil {
		t.Fatal("expected error for missing VALUE column, got nil")
	}
}

func TestPivot_MultiplePairs(t *testing.T) {
	entries := []Entry{
		{Key: "KEY", Value: "DB_HOST"},
		{Key: "VALUE", Value: "localhost"},
		{Key: "KEY", Value: "DB_PORT"},
		{Key: "VALUE", Value: "5432"},
	}
	out, err := Pivot(entries, DefaultPivotOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 pivoted entries, got %d", len(out))
	}
	if out[0].Key != "DB_HOST" || out[1].Key != "DB_PORT" {
		t.Errorf("unexpected keys: %v", out)
	}
}

func TestPivot_DefaultSeparatorFallback(t *testing.T) {
	opts := DefaultPivotOptions()
	opts.Separator = ""
	entries := []Entry{
		{Key: "KEY", Value: "REDIS_URL"},
		{Key: "VALUE", Value: "redis://localhost"},
	}
	out, err := Pivot(entries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 1 || out[0].Key != "REDIS_URL" {
		t.Errorf("unexpected result: %v", out)
	}
}
