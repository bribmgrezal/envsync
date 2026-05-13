package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/user/envsync/internal/envfile"
	"github.com/user/envsync/internal/export"
)

func runRedact(args []string) error {
	fs := flag.NewFlagSet("redact", flag.ContinueOnError)

	src := fs.String("src", ".env", "Source .env file to redact")
	out := fs.String("out", "", "Output file (defaults to stdout)")
	keys := fs.String("keys", "", "Comma-separated list of additional keys to redact")
	patterns := fs.String("patterns", "", "Comma-separated patterns to match against key names")
	placeholder := fs.String("placeholder", "***REDACTED***", "Replacement text for redacted values")
	format := fs.String("format", "dotenv", "Output format: dotenv, export, json")

	if err := fs.Parse(args); err != nil {
		return err
	}

	entries, err := envfile.Parse(*src)
	if err != nil {
		return fmt.Errorf("parse %q: %w", *src, err)
	}

	opts := envfile.DefaultRedactOptions()
	opts.Placeholder = *placeholder

	if *keys != "" {
		for _, k := range strings.Split(*keys, ",") {
			if k = strings.TrimSpace(k); k != "" {
				opts.Keys = append(opts.Keys, k)
			}
		}
	}

	if *patterns != "" {
		for _, p := range strings.Split(*patterns, ",") {
			if p = strings.TrimSpace(p); p != "" {
				opts.Patterns = append(opts.Patterns, p)
			}
		}
	}

	redacted := envfile.Redact(entries, opts)

	w := os.Stdout
	if *out != "" {
		f, err := os.Create(*out)
		if err != nil {
			return fmt.Errorf("create output file: %w", err)
		}
		defer f.Close()
		w = f
	}

	exportOpts := export.DefaultOptions()
	exportOpts.Format = *format
	return export.Write(w, redacted, exportOpts)
}
