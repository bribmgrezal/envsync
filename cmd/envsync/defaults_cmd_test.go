package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempDefaultsEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatalf("create temp env: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp env: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestRunDefaults_FillsMissingKeys(t *testing.T) {
	src := writeTempDefaultsEnv(t, "APP_NAME=myapp\nPORT=8080\n")
	defs := writeTempDefaultsEnv(t, "PORT=3000\nTIMEOUT=30s\n")
	out := filepath.Join(t.TempDir(), "out.env")

	err := runDefaults([]string{"--src", src, "--defaults", defs, "--out", out})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(out)
	if !strings.Contains(string(data), "TIMEOUT=30s") {
		t.Errorf("expected TIMEOUT=30s in output, got:\n%s", data)
	}
	if !strings.Contains(string(data), "PORT=8080") {
		t.Errorf("expected PORT=8080 (not overwritten), got:\n%s", data)
	}
}

func TestRunDefaults_OverwriteFlag(t *testing.T) {
	src := writeTempDefaultsEnv(t, "PORT=8080\n")
	defs := writeTempDefaultsEnv(t, "PORT=3000\n")
	out := filepath.Join(t.TempDir(), "out.env")

	err := runDefaults([]string{"--src", src, "--defaults", defs, "--overwrite", "--out", out})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(out)
	if !strings.Contains(string(data), "PORT=3000") {
		t.Errorf("expected PORT=3000 after overwrite, got:\n%s", data)
	}
}

func TestRunDefaults_MissingSource_ReturnsError(t *testing.T) {
	err := runDefaults([]string{"--src", "/nonexistent.env", "--defaults", "/nonexistent.env"})
	if err == nil {
		t.Error("expected error for missing source file, got nil")
	}
}

func TestRunDefaults_MissingFlags_ReturnsError(t *testing.T) {
	err := runDefaults([]string{})
	if err == nil {
		t.Error("expected error when flags are missing, got nil")
	}
}
