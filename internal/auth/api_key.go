package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// GenerateKey creates a new random API key of the form "dlab_<64 hex chars>".
// The prefix makes leaked keys identifiable/greppable in logs — a common
// convention (cf. Stripe's "sk_live_" keys).
func GenerateKey() (string, error) {
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return "", fmt.Errorf("generating random key: %w", err)
	}
	return "dlab_" + hex.EncodeToString(raw), nil
}

// HashKey returns the SHA-256 hex digest of a plaintext API key. Only
// this hash is ever stored — the plaintext key is shown once, at
// creation time, and never persisted or logged anywhere.
func HashKey(key string) string {
	sum := sha256.Sum256([]byte(key))
	return hex.EncodeToString(sum[:])
}
