package validate_test

import (
	"testing"

	"github.com/user/envsync/internal/validate"
)

func TestValidate_ValidKeys(t *testing.T) {
	env := map[string]string{
		"APP_NAME":    "envsync",
		"_PRIVATE":    "value",
		"DB_HOST":     "localhost",
		"VERSION_2":   "2.0",
	}
	result := validate.Validate(env)
	if !result.Valid() {
		t.Errorf("expected no issues, got: %v", result.Issues)
	}
}

func TestValidate_InvalidKeyFormat(t *testing.T) {
	env := map[string]string{
		"123INVALID": "value",
	}
	result := validate.Validate(env)
	if result.Valid() {
		t.Error("expected issues for invalid key format, got none")
	}
	if len(result.Issues) != 1 {
		t.Errorf("expected 1 issue, got %d", len(result.Issues))
	}
}

func TestValidate_WhitespaceOnlyValue(t *testing.T) {
	env := map[string]string{
		"MY_KEY": "   ",
	}
	result := validate.Validate(env)
	if result.Valid() {
		t.Error("expected issues for whitespace-only value, got none")
	}
}

func TestValidate_EmptyValueIsAllowed(t *testing.T) {
	env := map[string]string{
		"OPTIONAL_KEY": "",
	}
	result := validate.Validate(env)
	if !result.Valid() {
		t.Errorf("expected empty string value to be valid, got issues: %v", result.Issues)
	}
}

func TestValidate_EmptyEnv(t *testing.T) {
	env := map[string]string{}
	result := validate.Validate(env)
	if !result.Valid() {
		t.Errorf("expected empty env map to be valid, got issues: %v", result.Issues)
	}
}

func TestValidateKeys_AllPresent(t *testing.T) {
	env := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
	}
	required := []string{"DB_HOST", "DB_PORT"}
	result := validate.ValidateKeys(env, required)
	if !result.Valid() {
		t.Errorf("expected no issues, got: %v", result.Issues)
	}
}

func TestValidateKeys_MissingRequired(t *testing.T) {
	env := map[string]string{
		"DB_HOST": "localhost",
	}
	required := []string{"DB_HOST", "DB_PORT", "DB_NAME"}
	result := validate.ValidateKeys(env, required)
	if result.Valid() {
		t.Error("expected issues for missing required keys")
	}
	if len(result.Issues) != 2 {
		t.Errorf("expected 2 issues, got %d", len(result.Issues))
	}
}

func TestValidateKeys_EmptyRequiredValue(t *testing.T) {
	env := map[string]string{
		"SECRET_KEY": "",
	}
	required := []string{"SECRET_KEY"}
	result := validate.ValidateKeys(env, required)
	if result.Valid() {
		t.Error("expected issue for required key with empty value")
	}
}

func TestValidateKeys_EmptyRequiredList(t *testing.T) {
	env := map[string]string{
		"DB_HOST": "localhost",
	}
	required := []string{}
	result := validate.ValidateKeys(env, required)
	if !result.Valid() {
		t.Errorf("expected no issues for empty required list, got: %v", result.Issues)
	}
}

func TestIssue_String(t *testing.T) {
	issue := validate.Issue{Key: "MY_KEY", Message: "required key is missing"}
	got := issue.String()
	want := "[MY_KEY] required key is missing"
	if got != want {
		t.Errorf("Issue.String() = %q, want %q", got, want)
	}
}
