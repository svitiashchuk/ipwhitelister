// Authentication logic
package auth

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

// TokenValidityDuration specifies how long a login token is valid.
const TokenValidityDuration = 15 * time.Minute

// GenerateToken generates a secure, random token for email verification.
func GenerateToken() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
