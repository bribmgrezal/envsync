package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/user/envsync/internal/envfile"
	"github.com/user/envsync/internal/export"
)

// runSort parses flags and sorts the given .env file, writing the result to stdout or a file.
func runSort(args []string) error {
	fs := flag.NewFlagSet("sort", flag.ContinueOnError)

	order := fs.String("order", "alpha", "Sort order: alpha, alpha-desc, group")
	output := fs.String("output", "", "Output file path (default: stdout)")
	format := fs.String("format", "dotenv", "Output format: dotenv, export, json")

	if err := fs.Parse(args); err != nil {
		return err
	}

	positional := fs.Args()
	if len(positional) < 1 {
		return fmt.Errorf("sort: source file required")
	}

	entries, err := envfile.Parse(positional[0])
	if err != nil {
		return fmt.Errorf("sort: parse %q: %w", positional[0], err)
	}

	opts := envfile.DefaultSortOptions()
	switch *order {
	case "alpha":
		opts.Order = envfile.SortAlpha
	case "alpha-desc":
		opts.Order = envfile.SortAlphaDesc
	case "group":
		opts.Order = envfile.SortByGroup
	default:
		return fmt.Errorf("sort: unknown order %q (use alpha, alpha-desc, group)", *order)
	}

	sorted := envfile.Sort(entries, opts)

	w := os.Stdout
	if *output != "" {
		f, err := os.Create(*output)
		if err != nil {
			return fmt.Errorf("sort: create output file: %w", err)
		}
		defer f.Close()
		w = f
	}

	exportOpts := export.DefaultOptions()
	exportOpts.Format = export.FormatFromString(*format)

	if err := export.Write(w, sorted, exportOpts); err != nil {
		return fmt.Errorf("sort: write output: %w", err)
	}

	return nil
}
