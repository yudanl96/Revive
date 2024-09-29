package util

import (
	"math/rand"
)

func init() {
}

func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, length)

	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

func RandomStrLen() string {
	return RandomString(int(rand.Intn(10)) + 5)
}
