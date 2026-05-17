package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempMaskEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatalf("write temp env: %v", err)
	}
	return p
}

func TestRunMask_MasksSensitiveKeys(t *testing.T) {
	src := writeTempMaskEnv(t, "APP_NAME=envsync\nDB_PASSWORD=secret\nPORT=8080\n")
	out := filepath.Join(t.TempDir(), "out.env")

	if err := runMask([]string{"-src", src, "-out", out}); err != nil {
		t.Fatalf("runMask: %v", err)
	}

	data, _ := os.ReadFile(out)
	body := string(data)

	if strings.Contains(body, "secret") {
		t.Error("expected DB_PASSWORD to be masked")
	}
	if !strings.Contains(body, "envsync") {
		t.Error("expected APP_NAME value to be preserved")
	}
}

func TestRunMask_CustomPlaceholder(t *testing.T) {
	src := writeTempMaskEnv(t, "API_KEY=abc123\n")
	out := filepath.Join(t.TempDir(), "out.env")

	if err := runMask([]string{"-src", src, "-out", out, "-placeholder", "[HIDDEN]"}); err != nil {
		t.Fatalf("runMask: %v", err)
	}

	data, _ := os.ReadFile(out)
	if !strings.Contains(string(data), "[HIDDEN]") {
		t.Error("expected custom placeholder [HIDDEN] in output")
	}
}

func TestRunMask_ExplicitKeys(t *testing.T) {
	src := writeTempMaskEnv(t, "PORT=8080\nMY_CUSTOM=value\n")
	out := filepath.Join(t.TempDir(), "out.env")

	if err := runMask([]string{"-src", src, "-out", out, "-keys", "MY_CUSTOM"}); err != nil {
		t.Fatalf("runMask: %v", err)
	}

	data, _ := os.ReadFile(out)
	if strings.Contains(string(data), "value") {
		t.Error("expected MY_CUSTOM to be masked via explicit keys")
	}
}

func TestRunMask_MissingSource(t *testing.T) {
	err := runMask([]string{"-src", "/nonexistent/.env"})
	if err == nil {
		t.Error("expected error for missing source file")
	}
}
