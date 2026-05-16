// Package envfile provides utilities for parsing, transforming, and managing
// .env files.
//
// # Encrypt / Decrypt
//
// The Encrypt and Decrypt functions provide AES-256-GCM symmetric encryption
// for sensitive entry values within a slice of Entry objects.
//
// Only entries whose keys match one of the configured KeyPatterns are
// encrypted; all other entries pass through unchanged. Encrypted values are
// prefixed with a configurable marker (default "enc:") so that callers and
// downstream tools can distinguish encrypted from plaintext values.
//
// A passphrase is required; the 256-bit AES key is derived from it using
// SHA-256. Each value is encrypted with a fresh random nonce, so encrypting
// the same value twice produces different ciphertext.
//
// # Usage
//
//	opts := envfile.DefaultEncryptOptions()
//	opts.Passphrase = os.Getenv("ENVSYNC_PASSPHRASE")
//
//	encrypted, err := envfile.Encrypt(entries, opts)
//	plain,     err := envfile.Decrypt(encrypted, opts)
//
// Idempotency: calling Encrypt on already-encrypted entries is safe; values
// that already carry the prefix are left untouched.
package envfile
