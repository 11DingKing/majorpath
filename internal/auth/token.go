package auth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func NewToken() (string, error) {
	b := make([]byte, 32)
	if _, e := rand.Read(b); e != nil {
		return "", fmt.Errorf("token entropy: %w", e)
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
func TokenPreview(token string) string {
	if len(token) <= 8 {
		return token
	}
	return token[:4] + "..." + token[len(token)-4:]
}
