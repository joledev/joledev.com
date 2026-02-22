package services

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateToken creates a cryptographically secure random token (32 bytes = 64 hex chars).
func GenerateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
