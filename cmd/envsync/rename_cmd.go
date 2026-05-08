package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/user/envsync/internal/envfile"
)

// runRename implements the `envsync rename` sub-command.
// It reads a .env file, renames keys according to a JSON mapping, and writes
// the result back to the destination path (or stdout when dest is "-").
//
// Usage:
//
//	envsync rename -src=.env -mapping='{"OLD_KEY":"NEW_KEY"}' [-dest=out.env] [-strict]
func runRename(args []string) error {
	fs := flag.NewFlagSet("rename", flag.ContinueOnError)

	src := fs.String("src", ".env", "source .env file")
	dest := fs.String("dest", "-", "destination file; use - for stdout")
	mappingJSON := fs.String("mapping", "{}", "JSON object mapping old keys to new keys")
	strict := fs.Bool("strict", false, "error when a mapped key is not found in source")

	if err := fs.Parse(args); err != nil {
		return err
	}

	// Parse the key mapping.
	var mapping map[string]string
	if err := json.Unmarshal([]byte(*mappingJSON), &mapping); err != nil {
		return fmt.Errorf("invalid -mapping JSON: %w", err)
	}

	// Load source entries.
	entries, err := envfile.Parse(*src)
	if err != nil {
		return fmt.Errorf("reading %s: %w", *src, err)
	}

	opts := envfile.DefaultRenameOptions()
	opts.Mapping = mapping
	opts.SkipMissing = !*strict

	renamed, err := envfile.Rename(entries, opts)
	if err != nil {
		return err
	}

	// Write output.
	var out *os.File
	if *dest == "-" {
		out = os.Stdout
	} else {
		out, err = os.Create(*dest)
		if err != nil {
			return fmt.Errorf("creating %s: %w", *dest, err)
		}
		defer out.Close()
	}

	for _, e := range renamed {
		fmt.Fprintf(out, "%s=%s\n", e.Key, e.Value)
	}
	return nil
}
