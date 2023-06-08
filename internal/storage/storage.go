package storage

import (
	"sync"
)

type DB struct {
	urls map[string]string
	mu   *sync.Mutex
	p    *producer
	c    *consumer
}

func (db *DB) AddURL(url string, shortURL string) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.urls[shortURL] = url
	event := Event{
		ID:          len(db.urls),
		ShortURL:    shortURL,
		OriginalURL: url,
	}
	err := db.p.WriteEvent(&event)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) GetURL(shortURL string) (string, bool) {
	db.mu.Lock()
	defer db.mu.Unlock()
	bigURL, ok := db.urls[shortURL]
	return bigURL, ok
}

func NewDB(path string) (*DB, error) {
	p, err := NewProducer(path)
	if err != nil {
		return nil, err
	}

	c, err := NewConsumer(path)
	if err != nil {
		return nil, err
	}

	urls, err := c.GetMap()
	if err != nil {
		return nil, err
	}

	return &DB{
		urls: urls,
		mu:   &sync.Mutex{},
		p:    p,
		c:    c,
	}, nil
}
