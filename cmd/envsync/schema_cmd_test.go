package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func writeTempSchema(t *testing.T, fields []rawSchemaField) string {
	t.Helper()
	data, err := json.Marshal(fields)
	if err != nil {
		t.Fatal(err)
	}
	f, err := os.CreateTemp(t.TempDir(), "schema-*.json")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.Write(data); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func writeTempSchemaEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestRunSchema_PassesWhenValid(t *testing.T) {
	envPath := writeTempSchemaEnv(t, "APP_NAME=envsync\nPORT=8080\n")
	schemaPath := writeTempSchema(t, []rawSchemaField{
		{Key: "APP_NAME", Required: true},
		{Key: "PORT", Required: true, Pattern: `^\d+$`},
	})
	err := runSchema([]string{"--schema", schemaPath, "--env", envPath})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestRunSchema_FailsOnMissingRequiredKey(t *testing.T) {
	envPath := writeTempSchemaEnv(t, "PORT=9090\n")
	schemaPath := writeTempSchema(t, []rawSchemaField{
		{Key: "APP_NAME", Required: true},
		{Key: "PORT", Required: true},
	})
	err := runSchema([]string{"--schema", schemaPath, "--env", envPath})
	if err == nil {
		t.Fatal("expected error for missing required key")
	}
}

func TestRunSchema_FailsOnPatternMismatch(t *testing.T) {
	envPath := writeTempSchemaEnv(t, "PORT=not-a-port\n")
	schemaPath := writeTempSchema(t, []rawSchemaField{
		{Key: "PORT", Required: true, Pattern: `^\d+$`},
	})
	err := runSchema([]string{"--schema", schemaPath, "--env", envPath})
	if err == nil {
		t.Fatal("expected error for pattern mismatch")
	}
}

func TestRunSchema_MissingSchemaFlag_ReturnsError(t *testing.T) {
	err := runSchema([]string{})
	if err == nil {
		t.Fatal("expected error when --schema not provided")
	}
}

func TestRunSchema_InvalidSchemaJSON_ReturnsError(t *testing.T) {
	envPath := writeTempSchemaEnv(t, "KEY=value\n")
	f, _ := os.CreateTemp(t.TempDir(), "bad-*.json")
	f.WriteString("not valid json")
	f.Close()
	err := runSchema([]string{"--schema", f.Name(), "--env", envPath})
	if err == nil {
		t.Fatal("expected error for invalid JSON schema")
	}
}
