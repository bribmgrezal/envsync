// Package mask provides utilities for masking sensitive values
// in .env file entries before display or output.
package mask

import "strings"

// DefaultSecretKeys contains common key patterns considered sensitive.
var DefaultSecretKeys = []string{
	"SECRET",
	"PASSWORD",
	"PASSWD",
	"TOKEN",
	"API_KEY",
	"APIKEY",
	"PRIVATE_KEY",
	"AUTH",
	"CREDENTIAL",
	"ACCESS_KEY",
	"SIGNING_KEY",
}

// Masker holds configuration for masking sensitive values.
type Masker struct {
	SecretKeys []string
	MaskChar   string
	RevealLen  int
}

// New returns a Masker with default settings.
func New() *Masker {
	return &Masker{
		SecretKeys: DefaultSecretKeys,
		MaskChar:   "*",
		RevealLen:  0,
	}
}

// IsSensitive reports whether the given key matches any known secret pattern.
func (m *Masker) IsSensitive(key string) bool {
	upper := strings.ToUpper(key)
	for _, pattern := range m.SecretKeys {
		if strings.Contains(upper, pattern) {
			return true
		}
	}
	return false
}

// Mask returns a masked version of value if the key is sensitive.
// If RevealLen > 0, that many trailing characters are left visible.
func (m *Masker) Mask(key, value string) string {
	if !m.IsSensitive(key) {
		return value
	}
	if len(value) == 0 {
		return value
	}
	reveal := m.RevealLen
	if reveal >= len(value) {
		reveal = 0
	}
	masked := strings.Repeat(m.MaskChar, len(value)-reveal)
	if reveal > 0 {
		masked += value[len(value)-reveal:]
	}
	return masked
}

// MaskMap returns a copy of the provided map with sensitive values masked.
func (m *Masker) MaskMap(env map[string]string) map[string]string {
	result := make(map[string]string, len(env))
	for k, v := range env {
		result[k] = m.Mask(k, v)
	}
	return result
}
