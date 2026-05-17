package envfile

import (
	"testing"
)

func maskEntries() []Entry {
	return []Entry{
		{Key: "APP_NAME", Value: "envsync"},
		{Key: "DB_PASSWORD", Value: "s3cr3t"},
		{Key: "API_KEY", Value: "abc123"},
		{Key: "PORT", Value: "8080"},
		{Key: "AUTH_TOKEN", Value: "tok-xyz"},
	}
}

func TestMaskValues_MasksSensitiveKeys(t *testing.T) {
	result := MaskValues(maskEntries(), DefaultMaskOptions())

	for _, e := range result {
		switch e.Key {
		case "DB_PASSWORD", "API_KEY", "AUTH_TOKEN":
			if e.Value != "****" {
				t.Errorf("expected %q to be masked, got %q", e.Key, e.Value)
			}
			if !e.Masked {
				t.Errorf("expected %q Masked=true", e.Key)
			}
		case "APP_NAME", "PORT":
			if e.Masked {
				t.Errorf("expected %q Masked=false", e.Key)
			}
		}
	}
}

func TestMaskValues_CustomPlaceholder(t *testing.T) {
	opts := DefaultMaskOptions()
	opts.Placeholder = "[REDACTED]"
	result := MaskValues(maskEntries(), opts)

	for _, e := range result {
		if e.Key == "DB_PASSWORD" && e.Value != "[REDACTED]" {
			t.Errorf("expected [REDACTED], got %q", e.Value)
		}
	}
}

func TestMaskValues_ExplicitKeys(t *testing.T) {
	opts := DefaultMaskOptions()
	opts.ExplicitKeys = []string{"PORT"}
	result := MaskValues(maskEntries(), opts)

	for _, e := range result {
		if e.Key == "PORT" && !e.Masked {
			t.Errorf("expected PORT to be masked via ExplicitKeys")
		}
	}
}

func TestMaskValues_DoesNotMutateOriginal(t *testing.T) {
	src := maskEntries()
	origVal := src[1].Value // DB_PASSWORD
	MaskValues(src, DefaultMaskOptions())
	if src[1].Value != origVal {
		t.Errorf("original slice was mutated")
	}
}

func TestMaskValues_EmptyEntries(t *testing.T) {
	result := MaskValues([]Entry{}, DefaultMaskOptions())
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d entries", len(result))
	}
}

func TestMaskValues_DefaultPlaceholder_WhenEmpty(t *testing.T) {
	opts := MaskOptions{Placeholder: "", Patterns: DefaultMaskOptions().Patterns}
	result := MaskValues(maskEntries(), opts)
	for _, e := range result {
		if e.Key == "API_KEY" && e.Value != "****" {
			t.Errorf("expected default placeholder ****, got %q", e.Value)
		}
	}
}
