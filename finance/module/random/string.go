package random

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func String(length int) string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	return StringWithChars(length, chars)
}

func StringUpper(length int) string {
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	return StringWithChars(length, chars)
}

func StringWithChars(length int, chars string) string {
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}

	return string(result)
}

func Number(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func Token() string {
	return String(70)
}
