package envfile

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Entry represents a single key-value pair from an .env file.
type Entry struct {
	Key     string
	Value   string
	Comment string
	Line    int
}

// EnvFile holds the parsed contents of an .env file.
type EnvFile struct {
	Path    string
	Entries []Entry
	Index   map[string]*Entry
}

// Parse reads and parses an .env file from the given path.
func Parse(path string) (*EnvFile, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening env file: %w", err)
	}
	defer f.Close()

	env := &EnvFile{
		Path:  path,
		Index: make(map[string]*Entry),
	}

	scanner := bufio.NewScanner(f)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		raw := scanner.Text()
		trimmed := strings.TrimSpace(raw)

		if trimmed == "" {
			continue
		}

		if strings.HasPrefix(trimmed, "#") {
			continue
		}

		key, value, comment, err := parseLine(trimmed)
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", lineNum, err)
		}

		entry := Entry{
			Key:     key,
			Value:   value,
			Comment: comment,
			Line:    lineNum,
		}
		env.Entries = append(env.Entries, entry)
		env.Index[key] = &env.Entries[len(env.Entries)-1]
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanning env file: %w", err)
	}

	return env, nil
}

// parseLine splits a raw line into key, value, and optional inline comment.
func parseLine(line string) (key, value, comment string, err error) {
	eqIdx := strings.IndexByte(line, '=')
	if eqIdx < 0 {
		return "", "", "", fmt.Errorf("invalid line (missing '='): %q", line)
	}

	key = strings.TrimSpace(line[:eqIdx])
	rest := line[eqIdx+1:]

	if hashIdx := strings.Index(rest, " #"); hashIdx >= 0 {
		comment = strings.TrimSpace(rest[hashIdx+1:])
		rest = rest[:hashIdx]
	}

	value = strings.Trim(strings.TrimSpace(rest), `"`)
	return key, value, comment, nil
}
