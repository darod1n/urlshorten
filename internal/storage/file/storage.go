package file

import (
	"context"
	"errors"
	"sync"

	"github.com/darod1n/urlshorten/internal/helpers"
)

type DB struct {
	urls map[string]string
	path string
	mu   *sync.Mutex
	p    *producer
	c    *consumer
}

func (db *DB) AddURL(ctx context.Context, url string) (string, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	shortURL := helpers.GenerateShortURL(6)
	db.urls[shortURL] = url

	event := event{
		ID:          len(db.urls),
		ShortURL:    shortURL,
		OriginalURL: url,
	}

	err := db.p.WriteEvent(&event)
	if err != nil {
		return "", err
	}
	return shortURL, nil
}

func (db *DB) GetURL(ctx context.Context, shortURL string) (string, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	bigURL, ok := db.urls[shortURL]
	if ok {
		return bigURL, nil
	}
	return "", errors.New("failed get url")
}

func (db *DB) PingContext(ctx context.Context) error {
	return nil
}

func (db *DB) Close() error {
	return nil
}

func NewDB(path string) (*DB, error) {

	p, err := newProducer(path)
	if err != nil {
		return nil, err
	}

	c, err := newConsumer(path)
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
		path: path,
		p:    p,
		c:    c,
	}, nil
}
