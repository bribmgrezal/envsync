package envfile

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

var defaultsBase = []Entry{
	{Key: "APP_NAME", Value: "myapp"},
	{Key: "LOG_LEVEL", Value: ""},
	{Key: "PORT", Value: "8080"},
}

var defaultsDefs = []Entry{
	{Key: "LOG_LEVEL", Value: "info"},
	{Key: "PORT", Value: "3000"},
	{Key: "TIMEOUT", Value: "30s"},
}

func TestDefaults_FillsEmptyValues(t *testing.T) {
	opts := DefaultDefaultsOptions()
	out, err := Defaults(defaultsBase, defaultsDefs, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := ToMap(out)
	if m["LOG_LEVEL"] != "info" {
		t.Errorf("expected LOG_LEVEL=info, got %q", m["LOG_LEVEL"])
	}
}

func TestDefaults_DoesNotOverwriteExistingValues(t *testing.T) {
	opts := DefaultDefaultsOptions()
	out, err := Defaults(defaultsBase, defaultsDefs, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := ToMap(out)
	if m["PORT"] != "8080" {
		t.Errorf("expected PORT=8080, got %q", m["PORT"])
	}
}

func TestDefaults_AppendsMissingKeys(t *testing.T) {
	opts := DefaultDefaultsOptions()
	out, err := Defaults(defaultsBase, defaultsDefs, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := ToMap(out)
	if m["TIMEOUT"] != "30s" {
		t.Errorf("expected TIMEOUT=30s, got %q", m["TIMEOUT"])
	}
}

func TestDefaults_Overwrite_ReplacesExistingValues(t *testing.T) {
	opts := DefaultDefaultsOptions()
	opts.Overwrite = true
	out, err := Defaults(defaultsBase, defaultsDefs, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := ToMap(out)
	if m["PORT"] != "3000" {
		t.Errorf("expected PORT=3000, got %q", m["PORT"])
	}
}

func TestDefaults_SkipEmpty_OmitsBlankDefaults(t *testing.T) {
	opts := DefaultDefaultsOptions()
	defs := []Entry{{Key: "NEW_KEY", Value: ""}}
	out, err := Defaults(defaultsBase, defs, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, e := range out {
		if e.Key == "NEW_KEY" {
			t.Error("NEW_KEY should not have been added when SkipEmpty is true")
		}
	}
}

func TestDefaults_DoesNotMutateOriginal(t *testing.T) {
	opts := DefaultDefaultsOptions()
	original := make([]Entry, len(defaultsBase))
	copy(original, defaultsBase)
	_, _ = Defaults(defaultsBase, defaultsDefs, opts)
	if diff := cmp.Diff(original, defaultsBase); diff != "" {
		t.Errorf("original mutated:\n%s", diff)
	}
}

func TestDefaults_NilDefaults_ReturnsError(t *testing.T) {
	_, err := Defaults(defaultsBase, nil, DefaultDefaultsOptions())
	if err == nil {
		t.Error("expected error for nil defaults, got nil")
	}
}
