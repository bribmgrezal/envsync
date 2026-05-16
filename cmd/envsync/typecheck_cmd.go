package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/user/envsync/internal/envfile"
)

// runTypeCheck implements the `typecheck` sub-command.
// Usage: envsync typecheck -src=.env -rules=PORT:int,DEBUG:bool,API_URL:url
func runTypeCheck(args []string) error {
	fs := flag.NewFlagSet("typecheck", flag.ContinueOnError)
	src := fs.String("src", ".env", "source .env file to validate")
	rulesFlag := fs.String("rules", "", "comma-separated list of KEY:type[:required] rules")
	jsonOut := fs.Bool("json", false, "output violations as JSON")
	stopOnFirst := fs.Bool("stop-on-first", false, "stop after the first violation")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *rulesFlag == "" {
		return fmt.Errorf("typecheck: -rules is required")
	}

	rules, err := parseTypeRules(*rulesFlag)
	if err != nil {
		return fmt.Errorf("typecheck: %w", err)
	}

	entries, err := envfile.Parse(*src)
	if err != nil {
		return fmt.Errorf("typecheck: cannot read %s: %w", *src, err)
	}

	opts := envfile.DefaultTypeCheckOptions()
	opts.StopOnFirst = *stopOnFirst

	violations, _ := envfile.TypeCheck(entries, rules, opts)

	if *jsonOut {
		return json.NewEncoder(os.Stdout).Encode(violations)
	}

	if len(violations) == 0 {
		fmt.Println("typecheck: all values match declared types")
		return nil
	}

	for _, v := range violations {
		if v.Value == "" && strings.Contains(v.Reason, "missing") {
			fmt.Fprintf(os.Stderr, "  MISSING  %-24s expected type: %s\n", v.Key, v.Expected)
		} else {
			fmt.Fprintf(os.Stderr, "  INVALID  %-24s %s\n", v.Key, v.Reason)
		}
	}
	return fmt.Errorf("typecheck: %d violation(s) found", len(violations))
}

// parseTypeRules parses a comma-separated rule string such as
// "PORT:int,DEBUG:bool:required,API_URL:url".
func parseTypeRules(raw string) ([]envfile.TypeRule, error) {
	var rules []envfile.TypeRule
	for _, part := range strings.Split(raw, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		segments := strings.SplitN(part, ":", 3)
		if len(segments) < 2 {
			return nil, fmt.Errorf("invalid rule %q: expected KEY:type", part)
		}
		r := envfile.TypeRule{
			Key:  strings.TrimSpace(segments[0]),
			Type: strings.TrimSpace(segments[1]),
		}
		if len(segments) == 3 && strings.TrimSpace(segments[2]) == "required" {
			r.Required = true
		}
		rules = append(rules, r)
	}
	return rules, nil
}
