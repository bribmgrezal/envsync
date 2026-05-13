package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/user/envsync/internal/envfile"
	"github.com/user/envsync/internal/export"
)

// runInterpolate implements the `interpolate` sub-command.
// It reads a .env file, resolves variable references within values, and writes
// the result to stdout or an output file.
func runInterpolate(args []string) error {
	fs := flag.NewFlagSet("interpolate", flag.ContinueOnError)

	src := fs.String("src", "", "source .env file (required)")
	out := fs.String("out", "", "output file (default: stdout)")
	useOS := fs.Bool("use-os", true, "fall back to OS environment for unresolved references")
	failMissing := fs.Bool("fail-missing", false, "return an error when a variable reference cannot be resolved")
	fmt_ := fs.String("format", "dotenv", "output format: dotenv | export | json")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *src == "" {
		fs.Usage()
		return fmt.Errorf("--src is required")
	}

	entries, err := envfile.Parse(*src)
	if err != nil {
		return fmt.Errorf("parse %q: %w", *src, err)
	}

	opts := envfile.InterpolateOptions{
		UseOS:         *useOS,
		FailOnMissing: *failMissing,
	}

	resolved, err := envfile.Interpolate(entries, opts)
	if err != nil {
		return fmt.Errorf("interpolate: %w", err)
	}

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
	exportOpts.Format = *fmt_

	if err := export.Write(w, resolved, exportOpts); err != nil {
		return fmt.Errorf("write output: %w", err)
	}
	return nil
}
