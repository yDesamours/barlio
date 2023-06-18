package helper

import (
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func IsASet[T comparable](items []T) bool {
	var hashMap = make(map[T]bool)

	for _, item := range items {
		hashMap[item] = true
	}

	return len(hashMap) == len(items)
}

func StringIsNotEmpty[T ~string](s T) bool {
	return strings.TrimSpace(string(s)) != ""
}

func CompareHash[T ~string](src, target T) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(target), []byte(src)); err != nil {
		return false
	}
	return true
}

func HashPassword[T ~string](u T) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}
