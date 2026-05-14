package envfile

import (
	"regexp"
	"testing"
)

func TestValidateSchema_AllPresent_NoViolations(t *testing.T) {
	s := Schema{
		Fields: []SchemaField{
			{Key: "APP_NAME", Required: true},
			{Key: "PORT", Required: true},
		},
	}
	env := map[string]string{"APP_NAME": "envsync", "PORT": "8080"}
	got := ValidateSchema(s, env)
	if len(got) != 0 {
		t.Fatalf("expected no violations, got %v", got)
	}
}

func TestValidateSchema_MissingRequiredKey(t *testing.T) {
	s := Schema{
		Fields: []SchemaField{
			{Key: "DATABASE_URL", Required: true},
		},
	}
	env := map[string]string{}
	got := ValidateSchema(s, env)
	if len(got) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(got))
	}
	if got[0].Key != "DATABASE_URL" {
		t.Errorf("unexpected key: %s", got[0].Key)
	}
}

func TestValidateSchema_OptionalMissingKey_NoViolation(t *testing.T) {
	s := Schema{
		Fields: []SchemaField{
			{Key: "OPTIONAL_FLAG", Required: false},
		},
	}
	env := map[string]string{}
	got := ValidateSchema(s, env)
	if len(got) != 0 {
		t.Fatalf("expected no violations, got %v", got)
	}
}

func TestValidateSchema_PatternMismatch(t *testing.T) {
	s := Schema{
		Fields: []SchemaField{
			{Key: "PORT", Required: true, Pattern: regexp.MustCompile(`^\d+$`)},
		},
	}
	env := map[string]string{"PORT": "not-a-number"}
	got := ValidateSchema(s, env)
	if len(got) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(got))
	}
}

func TestValidateSchema_PatternMatch_NoViolation(t *testing.T) {
	s := Schema{
		Fields: []SchemaField{
			{Key: "PORT", Required: true, Pattern: regexp.MustCompile(`^\d+$`)},
		},
	}
	env := map[string]string{"PORT": "3000"}
	got := ValidateSchema(s, env)
	if len(got) != 0 {
		t.Fatalf("expected no violations, got %v", got)
	}
}

func TestValidateSchema_MultipleViolations(t *testing.T) {
	s := Schema{
		Fields: []SchemaField{
			{Key: "APP_NAME", Required: true},
			{Key: "PORT", Required: true, Pattern: regexp.MustCompile(`^\d+$`)},
		},
	}
	env := map[string]string{"PORT": "abc"}
	got := ValidateSchema(s, env)
	if len(got) != 2 {
		t.Fatalf("expected 2 violations, got %d: %v", len(got), got)
	}
}

func TestValidateSchema_WhitespaceOnlyValue_TreatedAsMissing(t *testing.T) {
	s := Schema{
		Fields: []SchemaField{
			{Key: "SECRET", Required: true},
		},
	}
	env := map[string]string{"SECRET": "   "}
	got := ValidateSchema(s, env)
	if len(got) != 1 {
		t.Fatalf("expected 1 violation for whitespace-only value, got %d", len(got))
	}
}
