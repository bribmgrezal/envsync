package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/user/envsync/internal/envfile"
	"github.com/user/envsync/internal/export"
)

// runTransform parses CLI flags for the transform sub-command and applies
// envfile.Transform to the source file, writing the result to stdout or a
// destination file.
func runTransform(args []string) error {
	fs := flag.NewFlagSet("transform", flag.ContinueOnError)

	src := fs.String("src", ".env", "Source .env file")
	dst := fs.String("dst", "", "Destination file (default: stdout)")
	uppercase := fs.Bool("uppercase", false, "Convert keys to uppercase")
	trimValues := fs.Bool("trim", true, "Trim whitespace from values")
	addPrefix := fs.String("add-prefix", "", "Prefix to add to every key")
	stripPrefix := fs.String("strip-prefix", "", "Prefix to remove from keys")

	if err := fs.Parse(args); err != nil {
		return err
	}

	entries, err := envfile.Parse(*src)
	if err != nil {
		return fmt.Errorf("transform: parse %q: %w", *src, err)
	}

	opts := envfile.TransformOptions{
		UppercaseKeys: *uppercase,
		TrimValues:    *trimValues,
		PrefixKeys:    *addPrefix,
		StripPrefix:   *stripPrefix,
	}

	transformed := envfile.Transform(entries, opts)

	writer := os.Stdout
	if *dst != "" {
		f, err := os.Create(*dst)
		if err != nil {
			return fmt.Errorf("transform: create %q: %w", *dst, err)
		}
		defer f.Close()
		writer = f
	}

	return export.Write(writer, transformed, export.DefaultOptions())
}
