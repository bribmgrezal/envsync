package envfile

import (
	"strings"
	"testing"
)

var encryptEntries = []Entry{
	{Key: "APP_NAME", Value: "myapp"},
	{Key: "DB_PASSWORD", Value: "s3cr3t"},
	{Key: "API_TOKEN", Value: "tok_abc123"},
	{Key: "LOG_LEVEL", Value: "info"},
	{Key: "AWS_SECRET_KEY", Value: "AKIAXXX"},
}

func TestEncrypt_EncryptsSensitiveValues(t *testing.T) {
	opts := DefaultEncryptOptions()
	opts.Passphrase = "test-passphrase"

	out, err := Encrypt(encryptEntries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, e := range out {
		if isSensitiveKey(e.Key, opts.KeyPatterns) {
			if !strings.HasPrefix(e.Value, opts.Prefix) {
				t.Errorf("key %q: expected encrypted prefix, got %q", e.Key, e.Value)
			}
		} else {
			if strings.HasPrefix(e.Value, opts.Prefix) {
				t.Errorf("key %q: safe key should not be encrypted", e.Key)
			}
		}
	}
}

func TestDecrypt_RestoresOriginalValues(t *testing.T) {
	opts := DefaultEncryptOptions()
	opts.Passphrase = "roundtrip-pass"

	encrypted, err := Encrypt(encryptEntries, opts)
	if err != nil {
		t.Fatalf("encrypt: %v", err)
	}
	decrypted, err := Decrypt(encrypted, opts)
	if err != nil {
		t.Fatalf("decrypt: %v", err)
	}
	for i, orig := range encryptEntries {
		if decrypted[i].Value != orig.Value {
			t.Errorf("key %q: want %q, got %q", orig.Key, orig.Value, decrypted[i].Value)
		}
	}
}

func TestEncrypt_AlreadyEncrypted_SkipsEntry(t *testing.T) {
	opts := DefaultEncryptOptions()
	opts.Passphrase = "idempotent-pass"

	first, _ := Encrypt(encryptEntries, opts)
	second, err := Encrypt(first, opts)
	if err != nil {
		t.Fatalf("second encrypt: %v", err)
	}
	for _, e := range second {
		if isSensitiveKey(e.Key, opts.KeyPatterns) {
			if strings.Count(e.Value, opts.Prefix) != 1 {
				t.Errorf("key %q: double-encrypted value detected", e.Key)
			}
		}
	}
}

func TestEncrypt_EmptyPassphrase_ReturnsError(t *testing.T) {
	opts := DefaultEncryptOptions()
	_, err := Encrypt(encryptEntries, opts)
	if err == nil {
		t.Fatal("expected error for empty passphrase")
	}
}

func TestDecrypt_WrongPassphrase_ReturnsError(t *testing.T) {
	encOpts := DefaultEncryptOptions()
	encOpts.Passphrase = "correct-pass"
	encrypted, _ := Encrypt(encryptEntries, encOpts)

	decOpts := DefaultEncryptOptions()
	decOpts.Passphrase = "wrong-pass"
	_, err := Decrypt(encrypted, decOpts)
	if err == nil {
		t.Fatal("expected error when decrypting with wrong passphrase")
	}
}

func TestEncrypt_DoesNotMutateOriginal(t *testing.T) {
	orig := make([]Entry, len(encryptEntries))
	copy(orig, encryptEntries)

	opts := DefaultEncryptOptions()
	opts.Passphrase = "mutation-check"
	_, _ = Encrypt(encryptEntries, opts)

	for i, e := range encryptEntries {
		if e.Value != orig[i].Value {
			t.Errorf("key %q: original mutated", e.Key)
		}
	}
}
