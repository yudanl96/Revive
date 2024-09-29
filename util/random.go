package util

import (
	"math/rand"
)

func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, length)

	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

func RandomShortStr() string {
	return RandomString(int(rand.Intn(10)) + 5)
}

func RandomLongStr() string {
	return RandomString(int(rand.Intn(50)) + 5)
}

func RandomInt(min, max int) int {
	return int(rand.Intn(max)) + min
}
