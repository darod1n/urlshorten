package storage

import (
	"sync"

	"go.uber.org/zap"
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

func NewDB(l *zap.SugaredLogger) (*DB, error) {
	path, err := getPath()
	if err != nil {
		return nil, err
	}

	p, errProducer := NewProducer(path)
	if errProducer != nil {
		return nil, errProducer
	}

	c, errConsumer := NewConsumer(path, l)
	if errConsumer != nil {
		return nil, errConsumer
	}

	return &DB{
		urls: c.GetMap(),
		mu:   &sync.Mutex{},
		p:    p,
		c:    c,
	}, nil
}
