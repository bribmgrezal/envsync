package export

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/user/envsync/internal/envfile"
)

// Format represents the output format for exported env data.
type Format string

const (
	FormatDotEnv Format = "dotenv"
	FormatExport Format = "export"
	FormatJSON   Format = "json"
)

// Options controls export behaviour.
type Options struct {
	Format  Format
	Sorted  bool
	Masked  bool
}

// DefaultOptions returns sensible export defaults.
func DefaultOptions() Options {
	return Options{
		Format: FormatDotEnv,
		Sorted: true,
		Masked: false,
	}
}

// Write serialises the given env entries to w in the requested format.
func Write(w io.Writer, entries []envfile.Entry, opts Options) error {
	if opts.Sorted {
		sorted := make([]envfile.Entry, len(entries))
		copy(sorted, entries)
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].Key < sorted[j].Key
		})
		entries = sorted
	}

	switch opts.Format {
	case FormatDotEnv:
		return writeDotEnv(w, entries, opts.Masked)
	case FormatExport:
		return writeExport(w, entries, opts.Masked)
	case FormatJSON:
		return writeJSON(w, entries, opts.Masked)
	default:
		return fmt.Errorf("export: unknown format %q", opts.Format)
	}
}

func writeDotEnv(w io.Writer, entries []envfile.Entry, masked bool) error {
	for _, e := range entries {
		v := value(e, masked)
		if _, err := fmt.Fprintf(w, "%s=%s\n", e.Key, quoteIfNeeded(v)); err != nil {
			return err
		}
	}
	return nil
}

func writeExport(w io.Writer, entries []envfile.Entry, masked bool) error {
	for _, e := range entries {
		v := value(e, masked)
		if _, err := fmt.Fprintf(w, "export %s=%s\n", e.Key, quoteIfNeeded(v)); err != nil {
			return err
		}
	}
	return nil
}

func writeJSON(w io.Writer, entries []envfile.Entry, masked bool) error {
	if _, err := fmt.Fprint(w, "{\n"); err != nil {
		return err
	}
	for i, e := range entries {
		v := value(e, masked)
		comma := ","
		if i == len(entries)-1 {
			comma = ""
		}
		if _, err := fmt.Fprintf(w, "  %q: %q%s\n", e.Key, v, comma); err != nil {
			return err
		}
	}
	_, err := fmt.Fprint(w, "}\n")
	return err
}

func value(e envfile.Entry, masked bool) string {
	if masked && e.Masked {
		return "***"
	}
	return e.Value
}

func quoteIfNeeded(v string) string {
	if strings.ContainsAny(v, " \t#") {
		return fmt.Sprintf("%q", v)
	}
	return v
}
