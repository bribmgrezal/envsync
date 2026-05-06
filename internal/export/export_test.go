package export_test

import (
	"strings"
	"testing"

	"github.com/user/envsync/internal/envfile"
	"github.com/user/envsync/internal/export"
)

func sampleEntries() []envfile.Entry {
	return []envfile.Entry{
		{Key: "APP_NAME", Value: "envsync", Masked: false},
		{Key: "DB_PASSWORD", Value: "s3cr3t", Masked: true},
		{Key: "PORT", Value: "8080", Masked: false},
	}
}

func TestWrite_DotEnvFormat(t *testing.T) {
	var buf strings.Builder
	opts := export.DefaultOptions()
	if err := export.Write(&buf, sampleEntries(), opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "APP_NAME=envsync") {
		t.Errorf("expected APP_NAME line, got:\n%s", out)
	}
	if !strings.Contains(out, "DB_PASSWORD=s3cr3t") {
		t.Errorf("expected unmasked DB_PASSWORD when masked=false, got:\n%s", out)
	}
}

func TestWrite_DotEnvFormat_Masked(t *testing.T) {
	var buf strings.Builder
	opts := export.DefaultOptions()
	opts.Masked = true
	if err := export.Write(&buf, sampleEntries(), opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "DB_PASSWORD=***") {
		t.Errorf("expected masked DB_PASSWORD, got:\n%s", out)
	}
	if strings.Contains(out, "s3cr3t") {
		t.Errorf("secret value must not appear in masked output")
	}
}

func TestWrite_ExportFormat(t *testing.T) {
	var buf strings.Builder
	opts := export.Options{Format: export.FormatExport, Sorted: true}
	if err := export.Write(&buf, sampleEntries(), opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, line := range strings.Split(strings.TrimSpace(out), "\n") {
		if !strings.HasPrefix(line, "export ") {
			t.Errorf("line missing 'export' prefix: %q", line)
		}
	}
}

func TestWrite_JSONFormat(t *testing.T) {
	var buf strings.Builder
	opts := export.Options{Format: export.FormatJSON, Sorted: true}
	if err := export.Write(&buf, sampleEntries(), opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.HasPrefix(out, "{") || !strings.HasSuffix(strings.TrimSpace(out), "}") {
		t.Errorf("expected JSON object, got:\n%s", out)
	}
}

func TestWrite_SortedOutput(t *testing.T) {
	var buf strings.Builder
	opts := export.Options{Format: export.FormatDotEnv, Sorted: true}
	if err := export.Write(&buf, sampleEntries(), opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	keys := make([]string, len(lines))
	for i, l := range lines {
		keys[i] = strings.SplitN(l, "=", 2)[0]
	}
	for i := 1; i < len(keys); i++ {
		if keys[i] < keys[i-1] {
			t.Errorf("output not sorted: %v", keys)
		}
	}
}

func TestWrite_UnknownFormat(t *testing.T) {
	var buf strings.Builder
	opts := export.Options{Format: "xml"}
	if err := export.Write(&buf, sampleEntries(), opts); err == nil {
		t.Error("expected error for unknown format, got nil")
	}
}
