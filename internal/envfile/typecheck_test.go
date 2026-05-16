package envfile

import (
	"testing"
)

func typecheckEntries() []Entry {
	return []Entry{
		{Key: "PORT", Value: "8080"},
		{Key: "RATE", Value: "3.14"},
		{Key: "DEBUG", Value: "true"},
		{Key: "API_URL", Value: "https://example.com"},
		{Key: "CONTACT", Value: "admin@example.com"},
		{Key: "APP_NAME", Value: "envsync"},
	}
}

func TestTypeCheck_NoViolations(t *testing.T) {
	rules := []TypeRule{
		{Key: "PORT", Type: "int"},
		{Key: "RATE", Type: "float"},
		{Key: "DEBUG", Type: "bool"},
		{Key: "API_URL", Type: "url"},
		{Key: "CONTACT", Type: "email"},
		{Key: "APP_NAME", Type: "string"},
	}
	v, _ := TypeCheck(typecheckEntries(), rules, DefaultTypeCheckOptions())
	if len(v) != 0 {
		t.Fatalf("expected no violations, got %d: %+v", len(v), v)
	}
}

func TestTypeCheck_InvalidInt(t *testing.T) {
	entries := []Entry{{Key: "PORT", Value: "not-a-number"}}
	rules := []TypeRule{{Key: "PORT", Type: "int"}}
	v, _ := TypeCheck(entries, rules, DefaultTypeCheckOptions())
	if len(v) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(v))
	}
	if v[0].Key != "PORT" {
		t.Errorf("unexpected key %q", v[0].Key)
	}
}

func TestTypeCheck_InvalidBool(t *testing.T) {
	entries := []Entry{{Key: "DEBUG", Value: "yes"}}
	rules := []TypeRule{{Key: "DEBUG", Type: "bool"}}
	v, _ := TypeCheck(entries, rules, DefaultTypeCheckOptions())
	if len(v) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(v))
	}
}

func TestTypeCheck_InvalidURL(t *testing.T) {
	entries := []Entry{{Key: "API_URL", Value: "ftp://bad.com"}}
	rules := []TypeRule{{Key: "API_URL", Type: "url"}}
	v, _ := TypeCheck(entries, rules, DefaultTypeCheckOptions())
	if len(v) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(v))
	}
}

func TestTypeCheck_InvalidEmail(t *testing.T) {
	entries := []Entry{{Key: "CONTACT", Value: "not-an-email"}}
	rules := []TypeRule{{Key: "CONTACT", Type: "email"}}
	v, _ := TypeCheck(entries, rules, DefaultTypeCheckOptions())
	if len(v) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(v))
	}
}

func TestTypeCheck_RequiredKeyMissing(t *testing.T) {
	entries := []Entry{{Key: "APP_NAME", Value: "envsync"}}
	rules := []TypeRule{
		{Key: "PORT", Type: "int", Required: true},
		{Key: "APP_NAME", Type: "string"},
	}
	v, _ := TypeCheck(entries, rules, DefaultTypeCheckOptions())
	if len(v) != 1 {
		t.Fatalf("expected 1 violation for missing required key, got %d", len(v))
	}
	if v[0].Key != "PORT" {
		t.Errorf("expected violation for PORT, got %q", v[0].Key)
	}
}

func TestTypeCheck_StopOnFirst(t *testing.T) {
	entries := []Entry{
		{Key: "PORT", Value: "bad"},
		{Key: "RATE", Value: "bad"},
	}
	rules := []TypeRule{
		{Key: "PORT", Type: "int"},
		{Key: "RATE", Type: "float"},
	}
	opts := TypeCheckOptions{StopOnFirst: true}
	v, _ := TypeCheck(entries, rules, opts)
	if len(v) != 1 {
		t.Fatalf("expected exactly 1 violation with StopOnFirst, got %d", len(v))
	}
}

func TestTypeCheck_UnknownKeysIgnored(t *testing.T) {
	entries := []Entry{
		{Key: "UNLISTED", Value: "anything"},
	}
	rules := []TypeRule{{Key: "PORT", Type: "int"}}
	v, _ := TypeCheck(entries, rules, DefaultTypeCheckOptions())
	if len(v) != 0 {
		t.Fatalf("expected no violations for unlisted key, got %d", len(v))
	}
}
