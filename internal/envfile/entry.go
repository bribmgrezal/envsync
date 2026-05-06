package envfile

// Entry represents a single key/value pair parsed from a .env file.
// The Masked field is set downstream by the mask package when the key
// matches a known sensitive pattern; it is not populated by the parser
// itself.
type Entry struct {
	// Key is the environment variable name, e.g. "DATABASE_URL".
	Key string

	// Value is the raw (unquoted) value associated with Key.
	Value string

	// Masked indicates whether the value should be hidden in reports
	// and exports. Set by the mask package after parsing.
	Masked bool

	// LineNo is the 1-based line number in the source file where this
	// entry was found. Useful for diagnostics and validation messages.
	LineNo int
}

// ToMap converts a slice of Entry values into a plain string map.
// Duplicate keys are resolved by keeping the last occurrence, matching
// the behaviour of most shells.
func ToMap(entries []Entry) map[string]string {
	m := make(map[string]string, len(entries))
	for _, e := range entries {
		m[e.Key] = e.Value
	}
	return m
}

// FromMap converts a string map into a slice of Entry values.
// The resulting entries have no LineNo or Masked information.
func FromMap(m map[string]string) []Entry {
	entries := make([]Entry, 0, len(m))
	for k, v := range m {
		entries = append(entries, Entry{Key: k, Value: v})
	}
	return entries
}
