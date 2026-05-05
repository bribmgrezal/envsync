package mask_test

import (
	"testing"

	"github.com/yourorg/envsync/internal/mask"
)

func TestIsSensitive_MatchesKnownPatterns(t *testing.T) {
	m := mask.New()
	sensitive := []string{"DB_PASSWORD", "API_KEY", "AUTH_TOKEN", "SECRET_KEY", "ACCESS_KEY_ID"}
	for _, key := range sensitive {
		if !m.IsSensitive(key) {
			t.Errorf("expected %q to be sensitive", key)
		}
	}
}

func TestIsSensitive_IgnoresSafeKeys(t *testing.T) {
	m := mask.New()
	safe := []string{"APP_ENV", "PORT", "LOG_LEVEL", "DATABASE_HOST"}
	for _, key := range safe {
		if m.IsSensitive(key) {
			t.Errorf("expected %q to NOT be sensitive", key)
		}
	}
}

func TestMask_SensitiveValue(t *testing.T) {
	m := mask.New()
	got := m.Mask("DB_PASSWORD", "supersecret")
	if got != "***********" {
		t.Errorf("expected all stars, got %q", got)
	}
}

func TestMask_SafeValue(t *testing.T) {
	m := mask.New()
	got := m.Mask("APP_ENV", "production")
	if got != "production" {
		t.Errorf("expected value unchanged, got %q", got)
	}
}

func TestMask_EmptyValue(t *testing.T) {
	m := mask.New()
	got := m.Mask("API_KEY", "")
	if got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}

func TestMask_RevealLen(t *testing.T) {
	m := mask.New()
	m.RevealLen = 4
	got := m.Mask("API_KEY", "abcdefgh")
	expected := "****efgh"
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestMaskMap(t *testing.T) {
	m := mask.New()
	env := map[string]string{
		"APP_ENV":     "staging",
		"DB_PASSWORD": "hunter2",
		"API_KEY":     "key123",
	}
	result := m.MaskMap(env)
	if result["APP_ENV"] != "staging" {
		t.Errorf("APP_ENV should be unchanged")
	}
	if result["DB_PASSWORD"] == "hunter2" {
		t.Errorf("DB_PASSWORD should be masked")
	}
	if result["API_KEY"] == "key123" {
		t.Errorf("API_KEY should be masked")
	}
}
