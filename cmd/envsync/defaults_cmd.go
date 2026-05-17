package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/owner/envsync/internal/envfile"
	"github.com/owner/envsync/internal/export"
)

func runDefaults(args []string) error {
	fs := flag.NewFlagSet("defaults", flag.ContinueOnError)

	src := fs.String("src", "", "primary .env file (required)")
	defaults := fs.String("defaults", "", "defaults .env file (required)")
	out := fs.String("out", "", "output file (stdout if omitted)")
	overwrite := fs.Bool("overwrite", false, "replace existing values with defaults")
	skipEmpty := fs.Bool("skip-empty", true, "ignore default entries with empty values")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *src == "" || *defaults == "" {
		return fmt.Errorf("defaults: --src and --defaults are required")
	}

	srcEntries, err := envfile.ParseFile(*src)
	if err != nil {
		return fmt.Errorf("defaults: reading src: %w", err)
	}

	defEntries, err := envfile.ParseFile(*defaults)
	if err != nil {
		return fmt.Errorf("defaults: reading defaults: %w", err)
	}

	opts := envfile.DefaultDefaultsOptions()
	opts.Overwrite = *overwrite
	opts.SkipEmpty = *skipEmpty

	result, err := envfile.Defaults(srcEntries, defEntries, opts)
	if err != nil {
		return fmt.Errorf("defaults: %w", err)
	}

	w := os.Stdout
	if *out != "" {
		f, err := os.Create(*out)
		if err != nil {
			return fmt.Errorf("defaults: creating output file: %w", err)
		}
		defer f.Close()
		w = f
	}

	return export.Write(w, result, export.DefaultOptions())
}
