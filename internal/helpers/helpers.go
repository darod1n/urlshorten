package helpers

import (
	"crypto/md5"
	"encoding/hex"
)

const (
	alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func GenerateShortURL(url string, len int) string {
	hash := md5.Sum([]byte(url))
	return hex.EncodeToString(hash[:])[:len]
}
