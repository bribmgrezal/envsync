package envfile

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strings"
)

// DefaultEncryptOptions returns sensible defaults for Encrypt.
func DefaultEncryptOptions() EncryptOptions {
	return EncryptOptions{
		Passphrase:  "",
		KeyPatterns: []string{"SECRET", "PASSWORD", "TOKEN", "KEY", "PRIVATE"},
		Prefix:      "enc:",
	}
}

// EncryptOptions controls which entries are encrypted and how.
type EncryptOptions struct {
	// Passphrase is used to derive the AES-256 key.
	Passphrase string
	// KeyPatterns lists substrings; matching keys will have their values encrypted.
	KeyPatterns []string
	// Prefix is prepended to encrypted values so they can be identified later.
	Prefix string
}

// Encrypt encrypts the values of sensitive entries in-place.
// Only entries whose keys match one of the KeyPatterns are encrypted.
// Already-encrypted values (those starting with Prefix) are left untouched.
func Encrypt(entries []Entry, opts EncryptOptions) ([]Entry, error) {
	if opts.Passphrase == "" {
		return nil, errors.New("encrypt: passphrase must not be empty")
	}
	key := deriveKey(opts.Passphrase)
	out := make([]Entry, len(entries))
	for i, e := range entries {
		if isSensitiveKey(e.Key, opts.KeyPatterns) && !strings.HasPrefix(e.Value, opts.Prefix) {
			enc, err := aesgcmEncrypt(key, e.Value)
			if err != nil {
				return nil, fmt.Errorf("encrypt: key %q: %w", e.Key, err)
			}
			e.Value = opts.Prefix + enc
		}
		out[i] = e
	}
	return out, nil
}

// Decrypt decrypts values that were previously encrypted by Encrypt.
func Decrypt(entries []Entry, opts EncryptOptions) ([]Entry, error) {
	if opts.Passphrase == "" {
		return nil, errors.New("decrypt: passphrase must not be empty")
	}
	key := deriveKey(opts.Passphrase)
	out := make([]Entry, len(entries))
	for i, e := range entries {
		if strings.HasPrefix(e.Value, opts.Prefix) {
			raw := strings.TrimPrefix(e.Value, opts.Prefix)
			plain, err := aesgcmDecrypt(key, raw)
			if err != nil {
				return nil, fmt.Errorf("decrypt: key %q: %w", e.Key, err)
			}
			e.Value = plain
		}
		out[i] = e
	}
	return out, nil
}

func deriveKey(passphrase string) []byte {
	h := sha256.Sum256([]byte(passphrase))
	return h[:]
}

func isSensitiveKey(key string, patterns []string) bool {
	upper := strings.ToUpper(key)
	for _, p := range patterns {
		if strings.Contains(upper, strings.ToUpper(p)) {
			return true
		}
	}
	return false
}

func aesgcmEncrypt(key []byte, plaintext string) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	sealed := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(sealed), nil
}

func aesgcmDecrypt(key []byte, encoded string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", fmt.Errorf("base64 decode: %w", err)
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	if len(data) < gcm.NonceSize() {
		return "", errors.New("ciphertext too short")
	}
	nonce, ciphertext := data[:gcm.NonceSize()], data[gcm.NonceSize():]
	plain, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("gcm open: %w", err)
	}
	return string(plain), nil
}
