package hash

import (
	"crypto/md5"
	"encoding/hex"
)

func Md5(input string) string {
	return hex.EncodeToString(Md5Bytes([]byte(input)))
}

func Md5Bytes(input []byte) []byte {
	hash := md5.New()
	hash.Write(input)
	return hash.Sum(nil)
}
