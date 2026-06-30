package hash

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
)

// GenerateNewRandomStr генерирует случайную строку длиной 32 байта
func GenerateNewRandomStr() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// ToHash хеширует строку
func ToHash(str string) string {
	h := sha256.Sum256([]byte(str))
	return base64.URLEncoding.EncodeToString(h[:])
}

func IsValidHash(hash, str string) bool {
	if hash == "" {
		return false
	}

	return hash == ToHash(str)
}
