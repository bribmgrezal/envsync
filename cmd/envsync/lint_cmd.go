package main

import (
	"fmt"
	"os"

	"github.com/user/envsync/internal/envfile"
	"github.com/user/envsync/internal/lint"
)

// runLint parses the given file and prints lint findings to stdout.
// It returns a non-zero exit code when any error-severity finding is present.
func runLint(path string) int {
	if path == "" {
		fmt.Fprintln(os.Stderr, "lint: no file specified")
		return 2
	}

	entries, err := envfile.Parse(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "lint: failed to parse %q: %v\n", path, err)
		return 2
	}

	findings := lint.Run(entries)
	if len(findings) == 0 {
		fmt.Printf("lint: %s — no issues found\n", path)
		return 0
	}

	hasError := false
	for _, f := range findings {
		fmt.Println(f)
		if f.Severity == lint.SeverityError {
			hasError = true
		}
	}

	if hasError {
		return 1
	}
	return 0
}
