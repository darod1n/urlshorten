package memory

import (
	"context"
	"errors"
	"sync"

	"github.com/darod1n/urlshorten/internal/helpers"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type DB struct {
	urls map[string]string
	mu   *sync.Mutex
}

func (db *DB) AddURL(ctx context.Context, url string) (string, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	shortURL := helpers.GenerateShortURL(6)
	db.urls[shortURL] = url
	return shortURL, nil
}

func (db *DB) GetURL(ctx context.Context, shortURL string) (string, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	bigURL, ok := db.urls[shortURL]
	if ok {
		return bigURL, nil
	}
	return "", errors.New("failed to get url")
}

func (db *DB) PingContext(ctx context.Context) error {
	return nil
}

func (db *DB) Close() error {
	return nil
}

func NewDB() (*DB, error) {
	return &DB{
		urls: make(map[string]string),
		mu:   &sync.Mutex{},
	}, nil
}
