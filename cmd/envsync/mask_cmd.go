package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/user/envsync/internal/envfile"
	"github.com/user/envsync/internal/export"
)

func runMask(args []string) error {
	fs := flag.NewFlagSet("mask", flag.ContinueOnError)

	src := fs.String("src", ".env", "Source .env file to mask")
	out := fs.String("out", "", "Output file (default: stdout)")
	placeholder := fs.String("placeholder", "****", "Replacement text for sensitive values")
	explicit := fs.String("keys", "", "Comma-separated list of additional keys to mask")
	format := fs.String("format", "dotenv", "Output format: dotenv, export, json")

	if err := fs.Parse(args); err != nil {
		return err
	}

	entries, err := envfile.Parse(*src)
	if err != nil {
		return fmt.Errorf("parse %q: %w", *src, err)
	}

	opts := envfile.DefaultMaskOptions()
	opts.Placeholder = *placeholder
	if *explicit != "" {
		for _, k := range strings.Split(*explicit, ",") {
			k = strings.TrimSpace(k)
			if k != "" {
				opts.ExplicitKeys = append(opts.ExplicitKeys, k)
			}
		}
	}

	masked := envfile.MaskValues(entries, opts)

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
	return export.Write(w, masked, exportOpts)
}
