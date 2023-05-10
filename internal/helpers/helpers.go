package helpers

import (
	"math/rand"
	"time"
)

func GenerateShortURL(lenToken int) string {
	var shortURL []byte
	rand.Seed(time.Now().UnixNano())
	alphabet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	for i := 0; i < lenToken; i++ {
		r := rand.Intn(len(alphabet))
		shortURL = append(shortURL, alphabet[r])
	}
	return string(shortURL)
}
