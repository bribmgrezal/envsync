package envfile

import (
	"testing"
)

func castEntries() []Entry {
	return []Entry{
		{Key: "PORT", Value: "08080"},
		{Key: "RATE", Value: "3.14000"},
		{Key: "ENABLED", Value: "TRUE"},
		{Key: "NAME", Value: "  hello  "},
		{Key: "RETRIES", Value: "not-a-number"},
	}
}

func TestCast_Int_NormalisesLeadingZero(t *testing.T) {
	opts := DefaultCastOptions()
	opts.Rules = []CastRule{{Key: "PORT", Type: CastInt}}
	out, results, err := Cast(castEntries(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Casted != "8080" {
		t.Errorf("expected '8080', got %q", results[0].Casted)
	}
	if out[0].Value != "8080" {
		t.Errorf("entry value not updated: %q", out[0].Value)
	}
}

func TestCast_Float_NormalisesTrailingZeros(t *testing.T) {
	opts := DefaultCastOptions()
	opts.Rules = []CastRule{{Key: "RATE", Type: CastFloat}}
	out, _, err := Cast(castEntries(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[1].Value != "3.14" {
		t.Errorf("expected '3.14', got %q", out[1].Value)
	}
}

func TestCast_Bool_LowercasesTrue(t *testing.T) {
	opts := DefaultCastOptions()
	opts.Rules = []CastRule{{Key: "ENABLED", Type: CastBool}}
	out, _, err := Cast(castEntries(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[2].Value != "true" {
		t.Errorf("expected 'true', got %q", out[2].Value)
	}
}

func TestCast_String_KeepsValueUnchanged(t *testing.T) {
	opts := DefaultCastOptions()
	opts.Rules = []CastRule{{Key: "NAME", Type: CastString}}
	out, _, err := Cast(castEntries(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[3].Value != "  hello  " {
		t.Errorf("expected original value, got %q", out[3].Value)
	}
}

func TestCast_Strict_ReturnsErrorOnBadInt(t *testing.T) {
	opts := DefaultCastOptions()
	opts.Strict = true
	opts.Rules = []CastRule{{Key: "RETRIES", Type: CastInt}}
	_, _, err := Cast(castEntries(), opts)
	if err == nil {
		t.Fatal("expected error for invalid int, got nil")
	}
}

func TestCast_NonStrict_KeepsOriginalOnFailure(t *testing.T) {
	opts := DefaultCastOptions()
	opts.Strict = false
	opts.Rules = []CastRule{{Key: "RETRIES", Type: CastInt}}
	out, results, err := Cast(castEntries(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Err == nil {
		t.Error("expected cast result to carry an error")
	}
	if out[4].Value != "not-a-number" {
		t.Errorf("expected original value preserved, got %q", out[4].Value)
	}
}

func TestCast_NoRules_ReturnsUnchanged(t *testing.T) {
	opts := DefaultCastOptions()
	original := castEntries()
	out, results, err := Cast(original, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected no results, got %d", len(results))
	}
	for i, e := range out {
		if e.Value != original[i].Value {
			t.Errorf("entry %d mutated: %q != %q", i, e.Value, original[i].Value)
		}
	}
}
