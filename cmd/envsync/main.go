// Command envsync diffs and syncs .env files across environments
// with secret masking support.
package main

import (
	"fmt"
	"os"

	"github.com/user/envsync/internal/diff"
	"github.com/user/envsync/internal/envfile"
	"github.com/user/envsync/internal/mask"
	"github.com/user/envsync/internal/report"
	"github.com/user/envsync/internal/sync"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	cfg, err := parseFlags(args)
	if err != nil {
		return err
	}

	source, err := envfile.Parse(cfg.sourceFile)
	if err != nil {
		return fmt.Errorf("parsing source %q: %w", cfg.sourceFile, err)
	}

	target, err := envfile.Parse(cfg.targetFile)
	if err != nil {
		return fmt.Errorf("parsing target %q: %w", cfg.targetFile, err)
	}

	results := diff.Compare(source, target)

	masker := mask.New(cfg.maskSecrets)
	renderer := report.New(os.Stdout, masker)
	renderer.Render(results)

	if cfg.apply {
		opts := sync.DefaultOptions()
		opts.RemoveExtra = cfg.removeExtra
		updated, changes := sync.Apply(target, results, opts)
		if changes > 0 {
			if err := envfile.Write(cfg.targetFile, updated); err != nil {
				return fmt.Errorf("writing target %q: %w", cfg.targetFile, err)
			}
			fmt.Fprintf(os.Stdout, "\nApplied %d change(s) to %s\n", changes, cfg.targetFile)
		} else {
			fmt.Fprintln(os.Stdout, "\nNo changes to apply.")
		}
	}

	return nil
}
