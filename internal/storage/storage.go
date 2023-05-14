package storage

import "sync"

type DB struct {
	urls map[string]string
	sync.Mutex
}

func (db *DB) AddURL(url string, shortURL string) {
	db.Lock()
	defer db.Unlock()
	db.urls[shortURL] = url
}

func (db *DB) GetURL(shortURL string) (string, bool) {
	db.Lock()
	defer db.Unlock()
	bigURL, ok := db.urls[shortURL]
	return bigURL, ok
}

func NewDB() *DB {
	return &DB{urls: map[string]string{}}
}
