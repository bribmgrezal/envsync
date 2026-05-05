package envfile

import (
	"os"
	"testing"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatalf("creating temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestParse_BasicKeyValues(t *testing.T) {
	path := writeTempEnv(t, "APP_ENV=production\nDB_HOST=localhost\nDB_PORT=5432\n")

	env, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(env.Entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(env.Entries))
	}

	if env.Index["APP_ENV"].Value != "production" {
		t.Errorf("expected APP_ENV=production, got %q", env.Index["APP_ENV"].Value)
	}
}

func TestParse_SkipsCommentsAndBlanks(t *testing.T) {
	path := writeTempEnv(t, "# this is a comment\n\nFOO=bar\n")

	env, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(env.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(env.Entries))
	}
}

func TestParse_QuotedValues(t *testing.T) {
	path := writeTempEnv(t, `SECRET="my secret value"\n`)

	env, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if v := env.Index["SECRET"].Value; v != "my secret value" {
		t.Errorf("expected unquoted value, got %q", v)
	}
}

func TestParse_InlineComment(t *testing.T) {
	path := writeTempEnv(t, "PORT=8080 # default port\n")

	env, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	entry := env.Index["PORT"]
	if entry.Value != "8080" {
		t.Errorf("expected value 8080, got %q", entry.Value)
	}
	if entry.Comment != "# default port" {
		t.Errorf("expected comment '# default port', got %q", entry.Comment)
	}
}

func TestParse_MissingEquals(t *testing.T) {
	path := writeTempEnv(t, "BADLINE\n")

	_, err := Parse(path)
	if err == nil {
		t.Fatal("expected error for malformed line, got nil")
	}
}

func TestParse_FileNotFound(t *testing.T) {
	_, err := Parse("/nonexistent/.env")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
