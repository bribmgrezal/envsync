package envfile

import (
	"testing"
)

func redactEntries() []Entry {
	return []Entry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "DB_PASSWORD", Value: "s3cr3t"},
		{Key: "API_KEY", Value: "abc123"},
		{Key: "DEBUG", Value: "true"},
		{Key: "AUTH_TOKEN", Value: "tok_xyz"},
		{Key: "PORT", Value: "8080"},
	}
}

func TestRedact_DefaultPatterns(t *testing.T) {
	result := Redact(redactEntries(), DefaultRedactOptions())

	expected := map[string]string{
		"APP_NAME": "myapp",
		"DB_PASSWORD": "***REDACTED***",
		"API_KEY":     "***REDACTED***",
		"DEBUG":       "true",
		"AUTH_TOKEN":  "***REDACTED***",
		"PORT":        "8080",
	}

	for _, e := range result {
		if got, ok := expected[e.Key]; ok {
			if e.Value != got {
				t.Errorf("key %q: got value %q, want %q", e.Key, e.Value, got)
			}
		}
	}
}

func TestRedact_ExplicitKeys(t *testing.T) {
	opts := RedactOptions{
		Keys:        []string{"APP_NAME", "PORT"},
		Placeholder: "[hidden]",
	}
	result := Redact(redactEntries(), opts)

	for _, e := range result {
		switch e.Key {
		case "APP_NAME", "PORT":
			if e.Value != "[hidden]" {
				t.Errorf("key %q should be redacted, got %q", e.Key, e.Value)
			}
			if !e.Masked {
				t.Errorf("key %q should have Masked=true", e.Key)
			}
		}
	}
}

func TestRedact_DoesNotMutateOriginal(t *testing.T) {
	orig := redactEntries()
	Redact(orig, DefaultRedactOptions())

	for _, e := range orig {
		if e.Key == "DB_PASSWORD" && e.Value != "s3cr3t" {
			t.Error("original entries were mutated")
		}
	}
}

func TestRedact_CustomPlaceholder(t *testing.T) {
	opts := DefaultRedactOptions()
	opts.Placeholder = "XXXXX"
	result := Redact(redactEntries(), opts)

	for _, e := range result {
		if e.Key == "API_KEY" && e.Value != "XXXXX" {
			t.Errorf("expected placeholder XXXXX, got %q", e.Value)
		}
	}
}

func TestRedact_EmptyEntries(t *testing.T) {
	result := Redact([]Entry{}, DefaultRedactOptions())
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d entries", len(result))
	}
}

func TestRedact_MaskedFlagSet(t *testing.T) {
	result := Redact(redactEntries(), DefaultRedactOptions())
	for _, e := range result {
		switch e.Key {
		case "DB_PASSWORD", "API_KEY", "AUTH_TOKEN":
			if !e.Masked {
				t.Errorf("key %q: expected Masked=true", e.Key)
			}
		case "APP_NAME", "DEBUG", "PORT":
			if e.Masked {
				t.Errorf("key %q: expected Masked=false", e.Key)
			}
		}
	}
}
