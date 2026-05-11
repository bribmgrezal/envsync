package main

import (
	"flag"
	"fmt"
	"os"
	"sort"

	"github.com/user/envsync/internal/envfile"
)

// runDedupe parses flags and removes duplicate keys from the given .env file,
// writing the result back to the source file (or stdout when --dry-run is set).
func runDedupe(args []string) error {
	fs := flag.NewFlagSet("dedupe", flag.ContinueOnError)
	keepFirst := fs.Bool("keep-first", false, "Keep the first occurrence of a duplicate key (default: last wins)")
	dryRun := fs.Bool("dry-run", false, "Print deduplicated output to stdout without writing to file")
	report := fs.Bool("report", false, "Print a summary of removed duplicate keys")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if fs.NArg() < 1 {
		return fmt.Errorf("usage: envsync dedupe [flags] <file>")
	}

	path := fs.Arg(0)

	entries, err := envfile.Parse(path)
	if err != nil {
		return fmt.Errorf("parse %q: %w", path, err)
	}

	opts := envfile.DefaultDedupeOptions()
	opts.KeepFirst = *keepFirst
	opts.ReportDuplicates = *report

	result := envfile.Dedupe(entries, opts)

	if *report && len(result.Duplicates) > 0 {
		sorted := append([]string(nil), result.Duplicates...)
		sort.Strings(sorted)
		fmt.Fprintf(os.Stderr, "dedupe: removed duplicates for keys: %v\n", sorted)
	}

	if *dryRun {
		for _, e := range result.Entries {
			if e.Value == "" {
				fmt.Fprintf(os.Stdout, "%s=\n", e.Key)
			} else {
				fmt.Fprintf(os.Stdout, "%s=%s\n", e.Key, e.Value)
			}
		}
		return nil
	}

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("open %q for writing: %w", path, err)
	}
	defer f.Close()

	for _, e := range result.Entries {
		if e.Value == "" {
			fmt.Fprintf(f, "%s=\n", e.Key)
		} else {
			fmt.Fprintf(f, "%s=%s\n", e.Key, e.Value)
		}
	}

	return nil
}
