package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

func HashPassword(password string) string {
	h := sha256.Sum256([]byte(password))
	return hex.EncodeToString(h[:])
}
func CheckPassword(hash, password string) bool {
	return strings.EqualFold(hash, HashPassword(password))
}
