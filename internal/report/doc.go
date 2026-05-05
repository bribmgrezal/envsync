// Package report provides rendering utilities for displaying the results
// of an environment file diff operation.
//
// It supports multiple output formats:
//
//   - FormatText: human-readable, line-oriented output suitable for terminals.
//   - FormatJSON: machine-readable JSON array of diff entries.
//
// Sensitive values are automatically masked using the mask.Masker before
// being written to the output stream, ensuring secrets are never exposed
// in plain text reports.
//
// Example usage:
//
//	r := report.New(mask.New(nil), report.FormatText)
//	_ = r.Render(os.Stdout, diffResults)
package report
