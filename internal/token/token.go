package token

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
)

func GenerateToken() (plainTextToken string, hashedToken []byte, err error) {
	randomBytes := make([]byte, 16)

	_, err = rand.Read(randomBytes)
	if err != nil {
		return
	}

	plainTextToken = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)

	hash := sha256.Sum256([]byte(plainTextToken))
	hashedToken = hash[:]
	return
}

func CompareToken(token string, hash []byte) bool {
	hashedToken := sha256.Sum256([]byte(token))
	return bytes.Compare(hashedToken[:], hash) == 0
}
