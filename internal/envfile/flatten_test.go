package envfile

import (
	"testing"
)

func flatEntries() []Entry {
	return []Entry{
		{Key: "DB__HOST", Value: "localhost"},
		{Key: "DB__PORT", Value: "5432"},
		{Key: "APP__SERVER__TIMEOUT", Value: "30"},
		{Key: "PLAIN", Value: "value"},
	}
}

func TestFlatten_DoubleSep_NormalisesToUnderscore(t *testing.T) {
	result := Flatten(flatEntries(), DefaultFlattenOptions())
	want := map[string]string{
		"DB_HOST":            "localhost",
		"DB_PORT":            "5432",
		"APP_SERVER_TIMEOUT": "30",
		"PLAIN":              "value",
	}
	for _, e := range result {
		exp, ok := want[e.Key]
		if !ok {
			t.Errorf("unexpected key %q", e.Key)
			continue
		}
		if e.Value != exp {
			t.Errorf("key %q: got value %q, want %q", e.Key, e.Value, exp)
		}
	}
}

func TestFlatten_DotSeparator(t *testing.T) {
	entries := []Entry{
		{Key: "db.host", Value: "localhost"},
		{Key: "db.port", Value: "5432"},
	}
	opts := DefaultFlattenOptions()
	result := Flatten(entries, opts)
	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
	if result[0].Key != "DB_HOST" {
		t.Errorf("expected DB_HOST, got %q", result[0].Key)
	}
}

func TestFlatten_WithPrefix(t *testing.T) {
	entries := []Entry{
		{Key: "HOST", Value: "localhost"},
	}
	opts := DefaultFlattenOptions()
	opts.Prefix = "APP"
	result := Flatten(entries, opts)
	if result[0].Key != "APP_HOST" {
		t.Errorf("expected APP_HOST, got %q", result[0].Key)
	}
}

func TestFlatten_UppercaseFalse_PreservesCase(t *testing.T) {
	entries := []Entry{
		{Key: "db__host", Value: "localhost"},
	}
	opts := DefaultFlattenOptions()
	opts.UppercaseKeys = false
	result := Flatten(entries, opts)
	if result[0].Key != "db_host" {
		t.Errorf("expected db_host, got %q", result[0].Key)
	}
}

func TestFlatten_NoSeparatorInKeys_PassesThrough(t *testing.T) {
	entries := []Entry{
		{Key: "SIMPLE", Value: "yes"},
	}
	result := Flatten(entries, DefaultFlattenOptions())
	if result[0].Key != "SIMPLE" {
		t.Errorf("expected SIMPLE, got %q", result[0].Key)
	}
}

func TestFlattenMap_NestedMap(t *testing.T) {
	m := map[string]interface{}{
		"db": map[string]interface{}{
			"host": "localhost",
			"port": 5432,
		},
		"debug": true,
	}
	result := FlattenMap(m, DefaultFlattenOptions())
	found := make(map[string]string)
	for _, e := range result {
		found[e.Key] = e.Value
	}
	if found["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q", found["DB_HOST"])
	}
	if found["DB_PORT"] != "5432" {
		t.Errorf("expected DB_PORT=5432, got %q", found["DB_PORT"])
	}
	if found["DEBUG"] != "true" {
		t.Errorf("expected DEBUG=true, got %q", found["DEBUG"])
	}
}
