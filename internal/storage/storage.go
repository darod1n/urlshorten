package storage

import "sync"

type DB struct {
	urls map[string]string
	mu   sync.Mutex
}

func (db *DB) AddURL(url string, shortURL string) {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.urls[shortURL] = url
}

func (db *DB) GetURL(shortURL string) (string, bool) {
	db.mu.Lock()
	defer db.mu.Unlock()
	bigURL, ok := db.urls[shortURL]
	return bigURL, ok
}

func NewDB() *DB {
	return &DB{urls: map[string]string{}}
}
