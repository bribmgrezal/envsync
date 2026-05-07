// Package diff provides functionality to compare two parsed .env files
// and report missing, extra, and changed keys between them.
package diff

// Status represents the type of difference found for a key.
type Status string

const (
	// StatusMissing indicates the key exists in base but not in target.
	StatusMissing Status = "missing"
	// StatusExtra indicates the key exists in target but not in base.
	StatusExtra Status = "extra"
	// StatusChanged indicates the key exists in both but with different values.
	StatusChanged Status = "changed"
	// StatusMatch indicates the key exists in both with the same value.
	StatusMatch Status = "match"
)

// Entry represents a single diff result for a key.
type Entry struct {
	Key       string
	Status    Status
	BaseValue string
	DestValue string
}

// Result holds the full diff between two env maps.
type Result struct {
	Entries []Entry
}

// HasDiff returns true if there are any non-matching entries.
func (r *Result) HasDiff() bool {
	for _, e := range r.Entries {
		if e.Status != StatusMatch {
			return true
		}
	}
	return false
}

// Summary returns counts of missing, extra, and changed entries.
func (r *Result) Summary() (missing, extra, changed int) {
	for _, e := range r.Entries {
		switch e.Status {
		case StatusMissing:
			missing++
		case StatusExtra:
			extra++
		case StatusChanged:
			changed++
		}
	}
	return
}

// Filter returns a new Result containing only entries with the given status.
func (r *Result) Filter(status Status) *Result {
	var entries []Entry
	for _, e := range r.Entries {
		if e.Status == status {
			entries = append(entries, e)
		}
	}
	return &Result{Entries: entries}
}

// Compare diffs base against dest, returning a Result with all entries.
// base is the reference environment (e.g. .env.example).
// dest is the target environment (e.g. .env.production).
func Compare(base, dest map[string]string) *Result {
	seen := make(map[string]bool)
	var entries []Entry

	for k, bv := range base {
		seen[k] = true
		if dv, ok := dest[k]; ok {
			if bv == dv {
				entries = append(entries, Entry{Key: k, Status: StatusMatch, BaseValue: bv, DestValue: dv})
			} else {
				entries = append(entries, Entry{Key: k, Status: StatusChanged, BaseValue: bv, DestValue: dv})
			}
		} else {
			entries = append(entries, Entry{Key: k, Status: StatusMissing, BaseValue: bv, DestValue: ""})
		}
	}

	for k, dv := range dest {
		if !seen[k] {
			entries = append(entries, Entry{Key: k, Status: StatusExtra, BaseValue: "", DestValue: dv})
		}
	}

	return &Result{Entries: entries}
}
