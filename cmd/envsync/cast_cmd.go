package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/yourorg/envsync/internal/envfile"
)

// runCast implements the `cast` sub-command.
// Usage: envsync cast --src=.env --dst=out.env KEY:TYPE [KEY:TYPE ...]
// TYPE is one of: string, int, float, bool
func runCast(args []string) error {
	fs := flag.NewFlagSet("cast", flag.ContinueOnError)
	src := fs.String("src", ".env", "source .env file")
	dst := fs.String("dst", "", "output file (default: overwrite src)")
	strict := fs.Bool("strict", true, "fail on coercion error")

	if err := fs.Parse(args); err != nil {
		return err
	}

	ruleArgs := fs.Args()
	if len(ruleArgs) == 0 {
		return fmt.Errorf("cast: at least one KEY:TYPE rule is required")
	}

	rules, err := parseCastRules(ruleArgs)
	if err != nil {
		return err
	}

	entries, err := envfile.ParseFile(*src)
	if err != nil {
		return fmt.Errorf("cast: read %s: %w", *src, err)
	}

	opts := envfile.DefaultCastOptions()
	opts.Rules = rules
	opts.Strict = *strict

	result, castResults, err := envfile.Cast(entries, opts)
	if err != nil {
		return fmt.Errorf("cast: %w", err)
	}

	for _, r := range castResults {
		if r.Err != nil {
			fmt.Fprintf(os.Stderr, "warning: %s: %v\n", r.Key, r.Err)
		} else {
			fmt.Fprintf(os.Stderr, "cast %s: %q -> %q (%s)\n", r.Key, r.Original, r.Casted, r.Type)
		}
	}

	outPath := *dst
	if outPath == "" {
		outPath = *src
	}

	f, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("cast: write %s: %w", outPath, err)
	}
	defer f.Close()

	for _, e := range result {
		fmt.Fprintf(f, "%s=%s\n", e.Key, e.Value)
	}
	return nil
}

func parseCastRules(args []string) ([]envfile.CastRule, error) {
	rules := make([]envfile.CastRule, 0, len(args))
	for _, arg := range args {
		parts := strings.SplitN(arg, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("cast: invalid rule %q — expected KEY:TYPE", arg)
		}
		rules = append(rules, envfile.CastRule{
			Key:  strings.TrimSpace(parts[0]),
			Type: envfile.CastType(strings.TrimSpace(parts[1])),
		})
	}
	return rules, nil
}
