package helper

import "strings"

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
