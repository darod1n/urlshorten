package helpers

import (
	"crypto/md5"
	"encoding/hex"
)

func GenerateShortURL(url string, len int) string {
	hash := md5.Sum([]byte(url))
	return hex.EncodeToString(hash[:])[:len]
}
