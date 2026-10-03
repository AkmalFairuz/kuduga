package hash

import (
	"crypto/sha256"
	"encoding/hex"
)

func Sha256(input string) string {
	sha := sha256.New()
	sha.Write([]byte(input))
	return hex.EncodeToString(sha.Sum(nil))
}

func Sha256Bytes(input []byte) []byte {
	sha := sha256.New()
	sha.Write(input)
	return sha.Sum(nil)
}
