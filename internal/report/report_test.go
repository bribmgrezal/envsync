package report_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envsync/internal/diff"
	"github.com/user/envsync/internal/mask"
	"github.com/user/envsync/internal/report"
)

func newRenderer(f report.Format) *report.Renderer {
	return report.New(mask.New(nil), f)
}

func TestRenderText_NoDiff(t *testing.T) {
	r := newRenderer(report.FormatText)
	var buf bytes.Buffer
	if err := r.Render(&buf, []diff.Result{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No differences") {
		t.Errorf("expected no-diff message, got: %q", buf.String())
	}
}

func TestRenderText_MissingKey(t *testing.T) {
	r := newRenderer(report.FormatText)
	results := []diff.Result{
		{Key: "DB_HOST", Status: diff.StatusMissing},
	}
	var buf bytes.Buffer
	_ = r.Render(&buf, results)
	if !strings.Contains(buf.String(), "MISSING") || !strings.Contains(buf.String(), "DB_HOST") {
		t.Errorf("expected MISSING line, got: %q", buf.String())
	}
}

func TestRenderText_ExtraKey(t *testing.T) {
	r := newRenderer(report.FormatText)
	results := []diff.Result{
		{Key: "EXTRA_VAR", Status: diff.StatusExtra, SourceValue: "hello"},
	}
	var buf bytes.Buffer
	_ = r.Render(&buf, results)
	if !strings.Contains(buf.String(), "EXTRA") {
		t.Errorf("expected EXTRA line, got: %q", buf.String())
	}
}

func TestRenderText_MasksSecrets(t *testing.T) {
	r := newRenderer(report.FormatText)
	results := []diff.Result{
		{Key: "SECRET_KEY", Status: diff.StatusChanged, SourceValue: "newsecret", TargetValue: "oldsecret"},
	}
	var buf bytes.Buffer
	_ = r.Render(&buf, results)
	output := buf.String()
	if strings.Contains(output, "newsecret") || strings.Contains(output, "oldsecret") {
		t.Errorf("secret values should be masked, got: %q", output)
	}
}

func TestRenderJSON_ContainsKey(t *testing.T) {
	r := newRenderer(report.FormatJSON)
	results := []diff.Result{
		{Key: "APP_ENV", Status: diff.StatusExtra, SourceValue: "production"},
	}
	var buf bytes.Buffer
	_ = r.Render(&buf, results)
	output := buf.String()
	if !strings.Contains(output, `"APP_ENV"`) {
		t.Errorf("expected key in JSON output, got: %q", output)
	}
	if !strings.Contains(output, `"extra"`) {
		t.Errorf("expected status in JSON output, got: %q", output)
	}
}

func TestRenderJSON_EmptyResults(t *testing.T) {
	r := newRenderer(report.FormatJSON)
	var buf bytes.Buffer
	_ = r.Render(&buf, []diff.Result{})
	output := strings.TrimSpace(buf.String())
	if !strings.HasPrefix(output, "[") || !strings.HasSuffix(output, "]") {
		t.Errorf("expected empty JSON array, got: %q", output)
	}
}
