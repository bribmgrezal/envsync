package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
)

// config holds parsed CLI configuration.
type config struct {
	sourceFile  string
	targetFile  string
	apply       bool
	removeExtra bool
	maskSecrets bool
}

func parseFlags(args []string) (*config, error) {
	fs := flag.NewFlagSet("envsync", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	var cfg config
	fs.StringVar(&cfg.sourceFile, "source", ".env", "source .env file (reference)")
	fs.StringVar(&cfg.targetFile, "target", ".env.local", "target .env file to diff/sync")
	fs.BoolVar(&cfg.apply, "apply", false, "apply changes from source to target")
	fs.BoolVar(&cfg.removeExtra, "remove-extra", false, "remove keys in target not present in source")
	fs.BoolVar(&cfg.maskSecrets, "mask", true, "mask sensitive values in report output")

	if err := fs.Parse(args); err != nil {
		printUsage(os.Stderr, fs)
		return nil, fmt.Errorf("parsing flags: %w", err)
	}

	if cfg.sourceFile == "" {
		return nil, errors.New("--source is required")
	}
	if cfg.targetFile == "" {
		return nil, errors.New("--target is required")
	}

	return &cfg, nil
}

func printUsage(w io.Writer, fs *flag.FlagSet) {
	fmt.Fprintln(w, "Usage: envsync [options]")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "Options:")
	fs.SetOutput(w)
	fs.PrintDefaults()
}
