package lint_test

import (
	"testing"

	"github.com/user/envsync/internal/envfile"
	"github.com/user/envsync/internal/lint"
)

func entries(pairs ...string) []envfile.Entry {
	var out []envfile.Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, envfile.Entry{Key: pairs[i], Value: pairs[i+1], LineNo: i/2 + 1})
	}
	return out
}

func TestRun_NoFindings(t *testing.T) {
	result := lint.Run(entries("APP_NAME", "envsync", "DEBUG", "false"))
	if len(result) != 0 {
		t.Fatalf("expected 0 findings, got %d: %v", len(result), result)
	}
}

func TestRun_DuplicateKey(t *testing.T) {
	e := []envfile.Entry{
		{Key: "FOO", Value: "bar", LineNo: 1},
		{Key: "FOO", Value: "baz", LineNo: 3},
	}
	findings := lint.Run(e)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Severity != lint.SeverityError {
		t.Errorf("expected error severity, got %s", findings[0].Severity)
	}
	if findings[0].Key != "FOO" {
		t.Errorf("expected key FOO, got %s", findings[0].Key)
	}
}

func TestRun_ShadowedOSKey(t *testing.T) {
	findings := lint.Run(entries("PATH", "/usr/local/bin"))
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Severity != lint.SeverityWarn {
		t.Errorf("expected warn severity, got %s", findings[0].Severity)
	}
}

func TestRun_WhitespaceOnlyValue(t *testing.T) {
	e := []envfile.Entry{{Key: "FOO", Value: "   ", LineNo: 1}}
	findings := lint.Run(e)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Severity != lint.SeverityWarn {
		t.Errorf("expected warn severity, got %s", findings[0].Severity)
	}
}

func TestRun_LongValue(t *testing.T) {
	long := make([]byte, 600)
	for i := range long {
		long[i] = 'x'
	}
	e := []envfile.Entry{{Key: "BIG", Value: string(long), LineNo: 1}}
	findings := lint.Run(e)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Severity != lint.SeverityWarn {
		t.Errorf("expected warn severity, got %s", findings[0].Severity)
	}
}

func TestRun_MultipleIssues(t *testing.T) {
	e := []envfile.Entry{
		{Key: "PATH", Value: "first", LineNo: 1},
		{Key: "PATH", Value: "second", LineNo: 2},
	}
	// PATH shadows OS var (warn) + duplicate (error)
	findings := lint.Run(e)
	if len(findings) != 3 {
		t.Fatalf("expected 3 findings (shadow + duplicate + shadow), got %d: %v", len(findings), findings)
	}
}

func TestFinding_String_WithLine(t *testing.T) {
	f := lint.Finding{Key: "FOO", Line: 5, Severity: lint.SeverityWarn, Message: "test msg"}
	s := f.String()
	if s == "" {
		t.Error("expected non-empty string")
	}
}
