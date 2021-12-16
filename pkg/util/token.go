package util

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateToken 生成 token
func GenerateToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}
