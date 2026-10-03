package hash

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
)

func HmacSha1(key, data string) string {
	return hex.EncodeToString(HmacSha1Bytes([]byte(key), []byte(data)))
}

func HmacSha1Bytes(key, data []byte) []byte {
	mac := hmac.New(sha1.New, key)
	mac.Write(data)
	return mac.Sum(nil)
}
