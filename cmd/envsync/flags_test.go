package main

import (
	"testing"
)

func TestParseFlags_Defaults(t *testing.T) {
	cfg, err := parseFlags([]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.sourceFile != ".env" {
		t.Errorf("expected source .env, got %q", cfg.sourceFile)
	}
	if cfg.targetFile != ".env.local" {
		t.Errorf("expected target .env.local, got %q", cfg.targetFile)
	}
	if cfg.apply {
		t.Error("expected apply=false by default")
	}
	if cfg.removeExtra {
		t.Error("expected removeExtra=false by default")
	}
	if !cfg.maskSecrets {
		t.Error("expected maskSecrets=true by default")
	}
}

func TestParseFlags_CustomValues(t *testing.T) {
	cfg, err := parseFlags([]string{
		"--source", "prod.env",
		"--target", "staging.env",
		"--apply",
		"--remove-extra",
		"--mask=false",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.sourceFile != "prod.env" {
		t.Errorf("expected source prod.env, got %q", cfg.sourceFile)
	}
	if cfg.targetFile != "staging.env" {
		t.Errorf("expected target staging.env, got %q", cfg.targetFile)
	}
	if !cfg.apply {
		t.Error("expected apply=true")
	}
	if !cfg.removeExtra {
		t.Error("expected removeExtra=true")
	}
	if cfg.maskSecrets {
		t.Error("expected maskSecrets=false")
	}
}

func TestParseFlags_UnknownFlag(t *testing.T) {
	_, err := parseFlags([]string{"--unknown-flag"})
	if err == nil {
		t.Fatal("expected error for unknown flag, got nil")
	}
}
