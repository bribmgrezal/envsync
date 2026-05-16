package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/user/envsync/internal/envfile"
	"github.com/user/envsync/internal/export"
)

// runClone implements the `clone` sub-command which duplicates one or more
// keys under new names, optionally dropping the originals.
func runClone(args []string) error {
	fs := flag.NewFlagSet("clone", flag.ContinueOnError)

	src := fs.String("src", "", "Source .env file (required)")
	out := fs.String("out", "", "Output file (defaults to stdout)")
	mapFlag := fs.String("map", "", "Comma-separated KEY=NEW_KEY pairs, e.g. DB_HOST=DATABASE_HOST,DB_PORT=PG_PORT")
	dropOriginal := fs.Bool("drop-original", false, "Remove the original key after cloning")
	skipMissing := fs.Bool("skip-missing", false, "Silently skip source keys not found in the file")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *src == "" {
		return fmt.Errorf("clone: --src is required")
	}
	if *mapFlag == "" {
		return fmt.Errorf("clone: --map is required")
	}

	entries, err := envfile.Parse(*src)
	if err != nil {
		return fmt.Errorf("clone: parse %q: %w", *src, err)
	}

	keyMap := make(map[string][]string)
	for _, pair := range strings.Split(*mapFlag, ",") {
		parts := strings.SplitN(strings.TrimSpace(pair), "=", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return fmt.Errorf("clone: invalid map entry %q (expected SRC=DST)", pair)
		}
		keyMap[parts[0]] = append(keyMap[parts[0]], parts[1])
	}

	opts := envfile.DefaultCloneOptions()
	opts.KeyMap = keyMap
	opts.KeepOriginal = !*dropOriginal
	opts.SkipMissing = *skipMissing

	result, err := envfile.Clone(entries, opts)
	if err != nil {
		return fmt.Errorf("clone: %w", err)
	}

	writer := os.Stdout
	if *out != "" {
		f, err := os.Create(*out)
		if err != nil {
			return fmt.Errorf("clone: create output %q: %w", *out, err)
		}
		defer f.Close()
		writer = f
	}

	return export.Write(writer, result, export.DefaultOptions())
}
