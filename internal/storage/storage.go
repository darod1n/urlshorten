package storage

import (
	"math/rand"
	"time"
)

type Storage struct {
	urls map[string]string
}

func (s *Storage) GenerateShortURL(lenToken int) string {
	var shortURL []byte
	rand.Seed(time.Now().UnixNano())
	alphabet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	for i := 0; i < lenToken; i++ {
		r := rand.Intn(len(alphabet))
		shortURL = append(shortURL, alphabet[r])
	}
	return string(shortURL)
}

func (s *Storage) AddToStorage(url string, shortURL string) {
	s.urls[shortURL] = url
}

func (s *Storage) GetURL(shortURL string) (string, bool) {
	bigURL, ok := s.urls[shortURL]
	return bigURL, ok
}

func NewStorage() *Storage {
	return &Storage{urls: map[string]string{}}
}
