package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempTransformEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatalf("writeTempTransformEnv: %v", err)
	}
	return p
}

func TestRunTransform_Uppercase(t *testing.T) {
	src := writeTempTransformEnv(t, "db_host=localhost\ndb_port=5432\n")
	dst := filepath.Join(t.TempDir(), "out.env")

	if err := runTransform([]string{"-src", src, "-dst", dst, "-uppercase"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(dst)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}
	out := string(data)
	if !strings.Contains(out, "DB_HOST") {
		t.Errorf("expected DB_HOST in output, got:\n%s", out)
	}
	if !strings.Contains(out, "DB_PORT") {
		t.Errorf("expected DB_PORT in output, got:\n%s", out)
	}
}

func TestRunTransform_AddPrefix(t *testing.T) {
	src := writeTempTransformEnv(t, "HOST=localhost\nPORT=5432\n")
	dst := filepath.Join(t.TempDir(), "out.env")

	if err := runTransform([]string{"-src", src, "-dst", dst, "-add-prefix", "APP_"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(dst)
	out := string(data)
	if !strings.Contains(out, "APP_HOST") {
		t.Errorf("expected APP_HOST in output, got:\n%s", out)
	}
}

func TestRunTransform_MissingSource(t *testing.T) {
	err := runTransform([]string{"-src", "/nonexistent/.env"})
	if err == nil {
		t.Error("expected error for missing source file")
	}
}
