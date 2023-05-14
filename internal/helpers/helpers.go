package helpers

import (
	"math/rand"
)

const (
	alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func GenerateShortURL(lenToken int) string {
	shortURL := make([]byte, 0, lenToken)
	for i := 0; i < lenToken; i++ {
		r := rand.Intn(len(alphabet))
		shortURL = append(shortURL, alphabet[r])
	}
	return string(shortURL)
}
