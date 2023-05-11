package helpers

import (
	"math/rand"
)

func GenerateShortURL(lenToken int) string {
	shortURL := make([]byte, 0, lenToken)
	alphabet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	for i := 0; i < lenToken; i++ {
		r := rand.Intn(len(alphabet))
		shortURL = append(shortURL, alphabet[r])
	}
	return string(shortURL)
}
