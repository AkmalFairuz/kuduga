package hash

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func HmacSha256(data, key string) string {
	return hex.EncodeToString(HmacSha256Bytes([]byte(data), []byte(key)))
}

func HmacSha256Bytes(data, key []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write(data)
	return mac.Sum(nil)
}
