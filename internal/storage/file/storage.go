package file

import (
	"context"
	"fmt"
	"net/url"
	"sync"

	"github.com/darod1n/urlshorten/internal/models"
	"github.com/darod1n/urlshorten/internal/storage/memory"
)

type DB struct {
	memory *memory.DB
	path   string
	uuid   int
	mu     *sync.Mutex
	p      *producer
	c      *consumer
}

func (db *DB) AddURL(ctx context.Context, url string) (string, error) {
	shortURL, err := db.memory.AddURL(ctx, url)
	if err != nil {
		return "", fmt.Errorf("failed to add url to memory storage: %v", err)
	}

	event := event{
		ID:          db.uuid + 1,
		ShortURL:    shortURL,
		OriginalURL: url,
	}

	if err := db.p.WriteEvent(&event); err != nil {
		return "", err
	}
	return shortURL, nil
}

func (db *DB) GetURL(ctx context.Context, shortURL string) (string, bool, error) {
	return db.memory.GetURL(ctx, shortURL)
}

func (db *DB) PingContext(ctx context.Context) error {
	return nil
}

func (db *DB) Close() {
}

func (db *DB) GetUserURLS(ctx context.Context, host string) ([]models.UserURLS, error) {
	return nil, nil
}

func (db *DB) Batch(ctx context.Context, host string, batch []models.BatchRequest) ([]models.BatchResponse, error) {
	var data []models.BatchResponse
	for _, val := range batch {
		shortURL, err := db.AddURL(ctx, val.OriginURL)
		if err != nil {
			return nil, fmt.Errorf("failed to add url: %v", err)
		}
		url, err := url.JoinPath(host, shortURL)
		if err != nil {
			return nil, fmt.Errorf("failed to join path: %v", err)
		}
		data = append(data, models.BatchResponse{CorrelationID: val.CorrelationID, ShortURL: url})
	}
	return data, nil
}

func (db *DB) DeleteUserURLS(ctx context.Context, userID string, urls []string) error {
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
	memory, err := memory.NewDB(urls)
	if err != nil {
		return nil, err
	}

	uuid := len(urls)

	return &DB{
		memory: memory,
		mu:     &sync.Mutex{},
		uuid:   uuid,
		path:   path,
		p:      p,
		c:      c,
	}, nil
}
