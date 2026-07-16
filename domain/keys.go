package domain

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

const KeyPrefix = "lnk_"

func GenerateAPIKey() (plaintext, hash, prefix string, err error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", "", "", err
	}
	plaintext = KeyPrefix + base64.RawURLEncoding.EncodeToString(b)
	hash = HashKey(plaintext)
	prefix = plaintext[:12]
	return plaintext, hash, prefix, nil
}

func HashKey(plaintext string) string {
	h := sha256.Sum256([]byte(plaintext))
	return hex.EncodeToString(h[:])
}
