// Package lint provides heuristic checks for .env file quality,
// flagging potential issues such as duplicate keys, suspicious values,
// overly long lines, and keys that shadow common OS variables.
package lint

import (
	"fmt"
	"strings"

	"github.com/user/envsync/internal/envfile"
)

// Severity indicates how serious a lint finding is.
type Severity string

const (
	SeverityWarn  Severity = "warn"
	SeverityError Severity = "error"
)

// Finding represents a single lint issue found in an env file.
type Finding struct {
	Key      string
	Line     int
	Severity Severity
	Message  string
}

func (f Finding) String() string {
	if f.Line > 0 {
		return fmt.Sprintf("[%s] line %d (%s): %s", f.Severity, f.Line, f.Key, f.Message)
	}
	return fmt.Sprintf("[%s] (%s): %s", f.Severity, f.Key, f.Message)
}

// shadowedKeys are well-known OS-level variables that should not be
// overridden lightly inside application .env files.
var shadowedKeys = map[string]bool{
	"PATH": true, "HOME": true, "USER": true, "SHELL": true,
	"LANG": true, "PWD": true, "TERM": true,
}

const maxLineLength = 512

// Run performs all lint checks on the provided entries and returns a
// (possibly empty) slice of Findings.
func Run(entries []envfile.Entry) []Finding {
	var findings []Finding
	seen := make(map[string]int) // key -> first line number

	for _, e := range entries {
		// Duplicate key check.
		if first, ok := seen[e.Key]; ok {
			findings = append(findings, Finding{
				Key:      e.Key,
				Line:     e.LineNo,
				Severity: SeverityError,
				Message:  fmt.Sprintf("duplicate key (first seen on line %d)", first),
			})
		} else {
			seen[e.Key] = e.LineNo
		}

		// Shadowed OS variable check.
		if shadowedKeys[strings.ToUpper(e.Key)] {
			findings = append(findings, Finding{
				Key:      e.Key,
				Line:     e.LineNo,
				Severity: SeverityWarn,
				Message:  "key shadows a well-known OS environment variable",
			})
		}

		// Overly long value check.
		if len(e.Value) > maxLineLength {
			findings = append(findings, Finding{
				Key:      e.Key,
				Line:     e.LineNo,
				Severity: SeverityWarn,
				Message:  fmt.Sprintf("value exceeds %d characters", maxLineLength),
			})
		}

		// Whitespace-only value check.
		if len(e.Value) > 0 && strings.TrimSpace(e.Value) == "" {
			findings = append(findings, Finding{
				Key:      e.Key,
				Line:     e.LineNo,
				Severity: SeverityWarn,
				Message:  "value contains only whitespace",
			})
		}
	}

	return findings
}
