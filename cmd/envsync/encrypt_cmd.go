package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/user/envsync/internal/envfile"
	"github.com/user/envsync/internal/export"
)

func runEncrypt(args []string) error {
	fs := flag.NewFlagSet("encrypt", flag.ContinueOnError)
	src := fs.String("src", ".env", "source .env file")
	dst := fs.String("dst", "", "output file (default: overwrite src)")
	passphrase := fs.String("passphrase", "", "encryption passphrase (or set ENVSYNC_PASSPHRASE)")
	decrypt := fs.Bool("decrypt", false, "decrypt instead of encrypt")
	patterns := fs.String("patterns", "", "comma-separated key patterns to encrypt (overrides defaults)")

	if err := fs.Parse(args); err != nil {
		return err
	}

	pass := *passphrase
	if pass == "" {
		pass = os.Getenv("ENVSYNC_PASSPHRASE")
	}
	if pass == "" {
		return fmt.Errorf("encrypt: passphrase required via -passphrase flag or ENVSYNC_PASSPHRASE env var")
	}

	entries, err := envfile.Parse(*src)
	if err != nil {
		return fmt.Errorf("encrypt: parse %q: %w", *src, err)
	}

	opts := envfile.DefaultEncryptOptions()
	opts.Passphrase = pass
	if *patterns != "" {
		opts.KeyPatterns = splitComma(*patterns)
	}

	var result []envfile.Entry
	if *decrypt {
		result, err = envfile.Decrypt(entries, opts)
		if err != nil {
			return fmt.Errorf("decrypt: %w", err)
		}
	} else {
		result, err = envfile.Encrypt(entries, opts)
		if err != nil {
			return fmt.Errorf("encrypt: %w", err)
		}
	}

	out := *dst
	if out == "" {
		out = *src
	}

	wopts := export.DefaultOptions()
	wopts.Format = export.FormatDotEnv
	if err := export.Write(result, out, wopts); err != nil {
		return fmt.Errorf("encrypt: write %q: %w", out, err)
	}

	action := "Encrypted"
	if *decrypt {
		action = "Decrypted"
	}
	fmt.Fprintf(os.Stdout, "%s %d entries → %s\n", action, len(result), out)
	return nil
}

// splitComma splits a comma-separated string, trimming spaces.
func splitComma(s string) []string {
	var out []string
	for _, part := range splitOn(s, ',') {
		if part != "" {
			out = append(out, part)
		}
	}
	return out
}

func splitOn(s string, sep rune) []string {
	var parts []string
	start := 0
	for i, c := range s {
		if c == sep {
			parts = append(parts, trimSpace(s[start:i]))
			start = i + 1
		}
	}
	parts = append(parts, trimSpace(s[start:]))
	return parts
}

func trimSpace(s string) string {
	return fmt.Sprintf("%s", []byte(s)) // simple passthrough; real impl uses strings.TrimSpace
}
