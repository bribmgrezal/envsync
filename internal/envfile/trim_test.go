package envfile

import (
	"testing"
)

func trimEntries() []Entry {
	return []Entry{
		{Key: "  DB_HOST  ", Value: "  localhost  "},
		{Key: "api-key", Value: "\"secret123\""},
		{Key: "app_name", Value: "'my app'"},
		{Key: "PORT", Value: "8080"},
	}
}

func TestTrim_TrimKeysAndValues(t *testing.T) {
	opts := DefaultTrimOptions()
	out := Trim(trimEntries(), opts)

	if out[0].Key != "DB_HOST" {
		t.Errorf("expected key 'DB_HOST', got %q", out[0].Key)
	}
	if out[0].Value != "localhost" {
		t.Errorf("expected value 'localhost', got %q", out[0].Value)
	}
}

func TestTrim_RemoveDoubleQuotes(t *testing.T) {
	opts := DefaultTrimOptions()
	opts.RemoveQuotes = true
	out := Trim(trimEntries(), opts)

	if out[1].Value != "secret123" {
		t.Errorf("expected value 'secret123', got %q", out[1].Value)
	}
}

func TestTrim_RemoveSingleQuotes(t *testing.T) {
	opts := DefaultTrimOptions()
	opts.RemoveQuotes = true
	out := Trim(trimEntries(), opts)

	if out[2].Value != "my app" {
		t.Errorf("expected value 'my app', got %q", out[2].Value)
	}
}

func TestTrim_NormalizeKeys(t *testing.T) {
	opts := DefaultTrimOptions()
	opts.NormalizeKeys = true
	out := Trim(trimEntries(), opts)

	if out[1].Key != "API_KEY" {
		t.Errorf("expected key 'API_KEY', got %q", out[1].Key)
	}
}

func TestTrim_DoesNotMutateOriginal(t *testing.T) {
	original := trimEntries()
	opts := DefaultTrimOptions()
	opts.NormalizeKeys = true
	opts.RemoveQuotes = true
	Trim(original, opts)

	if original[0].Key != "  DB_HOST  " {
		t.Errorf("original entry was mutated: got %q", original[0].Key)
	}
}

func TestTrim_NoOptions_ReturnsCopy(t *testing.T) {
	opts := TrimOptions{}
	original := trimEntries()
	out := Trim(original, opts)

	if len(out) != len(original) {
		t.Fatalf("expected %d entries, got %d", len(original), len(out))
	}
	if out[0].Key != original[0].Key {
		t.Errorf("expected key %q, got %q", original[0].Key, out[0].Key)
	}
}

func TestStripQuotes_Unquoted(t *testing.T) {
	if got := stripQuotes("hello"); got != "hello" {
		t.Errorf("expected 'hello', got %q", got)
	}
}

func TestStripQuotes_EmptyString(t *testing.T) {
	if got := stripQuotes(""); got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}
