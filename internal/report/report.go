package report

import (
	"fmt"
	"io"
	"strings"

	"github.com/user/envsync/internal/diff"
	"github.com/user/envsync/internal/mask"
)

// Format controls the output style of a report.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Renderer writes a diff result to an output stream.
type Renderer struct {
	masker *mask.Masker
	format Format
}

// New creates a new Renderer with the given masker and output format.
func New(m *mask.Masker, f Format) *Renderer {
	return &Renderer{masker: m, format: f}
}

// Render writes the diff results to w.
func (r *Renderer) Render(w io.Writer, results []diff.Result) error {
	if r.format == FormatJSON {
		return r.renderJSON(w, results)
	}
	return r.renderText(w, results)
}

func (r *Renderer) renderText(w io.Writer, results []diff.Result) error {
	if len(results) == 0 {
		_, err := fmt.Fprintln(w, "No differences found.")
		return err
	}
	for _, res := range results {
		var line string
		switch res.Status {
		case diff.StatusMatch:
			continue
		case diff.StatusMissing:
			line = fmt.Sprintf("- MISSING  %s", res.Key)
		case diff.StatusExtra:
			val := r.masker.Mask(res.Key, res.SourceValue)
			line = fmt.Sprintf("+ EXTRA    %s=%s", res.Key, val)
		case diff.StatusChanged:
			src := r.masker.Mask(res.Key, res.SourceValue)
			dst := r.masker.Mask(res.Key, res.TargetValue)
			line = fmt.Sprintf("~ CHANGED  %s: %s -> %s", res.Key, dst, src)
		}
		if _, err := fmt.Fprintln(w, line); err != nil {
			return err
		}
	}
	return nil
}

func (r *Renderer) renderJSON(w io.Writer, results []diff.Result) error {
	var sb strings.Builder
	sb.WriteString("[\n")
	first := true
	for _, res := range results {
		if res.Status == diff.StatusMatch {
			continue
		}
		if !first {
			sb.WriteString(",\n")
		}
		src := r.masker.Mask(res.Key, res.SourceValue)
		dst := r.masker.Mask(res.Key, res.TargetValue)
		sb.WriteString(fmt.Sprintf(
			`  {"key":%q,"status":%q,"source":%q,"target":%q}`,
			res.Key, res.Status, src, dst,
		))
		first = false
	}
	sb.WriteString("\n]")
	_, err := fmt.Fprintln(w, sb.String())
	return err
}
