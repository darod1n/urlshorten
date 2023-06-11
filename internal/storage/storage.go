package storage

import (
	"context"
	"database/sql"
	"sync"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type DB struct {
	urls map[string]string
	path string
	mu   *sync.Mutex
	p    *producer
	c    *consumer
	base *sql.DB
}

func (db *DB) AddURL(url string, shortURL string) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.urls[shortURL] = url

	if db.path == "" {
		return nil
	}

	event := event{
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

func (db *DB) PingContext(ctx context.Context) error {
	return db.base.PingContext(ctx)
}

func (db *DB) Close() error {
	return db.base.Close()
}

func NewDB(path, driverName, dataSourceName string) (*DB, error) {
	if path == "" {
		return &DB{
			urls: make(map[string]string),
			mu:   &sync.Mutex{},
			path: path,
		}, nil
	}

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

	base, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}

	return &DB{
		urls: urls,
		mu:   &sync.Mutex{},
		path: path,
		p:    p,
		c:    c,
		base: base,
	}, nil
}
