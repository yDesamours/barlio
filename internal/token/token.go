package token

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
)

func GenerateToken() (plainTextToken, hashedToken string, err error) {
	randomBytes := make([]byte, 16)

	_, err = rand.Read(randomBytes)
	if err != nil {
		return
	}

	plainTextToken = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)

	hash := sha256.Sum256([]byte(plainTextToken))
	hashedToken = string(hash[:])
	return
}
