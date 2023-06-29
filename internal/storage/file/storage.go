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
		return "", fmt.Errorf("failed to add url to memory storage: %w", err)
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

func (db *DB) GetURL(ctx context.Context, shortURL string) (string, error) {
	return db.memory.GetURL(ctx, shortURL)
}

func (db *DB) PingContext(ctx context.Context) error {
	return db.memory.PingContext(ctx)
}

func (db *DB) Close() {
	db.memory.Close()
}

func (db *DB) GetUserURLS(ctx context.Context, host string) ([]models.UserURLS, error) {
	return db.memory.GetUserURLS(ctx, host)
}

func (db *DB) Batch(ctx context.Context, host string, br []models.BatchRequest) ([]models.BatchResponse, error) {
	bResp := make([]models.BatchResponse, len(br))
	for _, val := range br {
		shortURL, err := db.AddURL(ctx, val.OriginURL)
		if err != nil {
			return nil, fmt.Errorf("failed to add url: %w", err)
		}
		url, err := url.JoinPath(host, shortURL)
		if err != nil {
			return nil, fmt.Errorf("failed to join path: %w", err)
		}
		bResp = append(bResp, models.BatchResponse{CorrelationID: val.CorrelationID, ShortURL: url})
	}
	return bResp, nil
}

func (db *DB) DeleteUserURLS(ctx context.Context, userID string, urls []string) error {
	return db.memory.DeleteUserURLS(ctx, userID, urls)
}

func NewDB(path string) (*DB, error) {

	p, err := newProducer(path)
	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %w", err)
	}

	c, err := newConsumer(path)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}

	urls, err := c.GetMap()
	if err != nil {
		return nil, err
	}
	memory, err := memory.NewDB(urls)
	if err != nil {
		return nil, fmt.Errorf("failed to create memory db: %w", err)
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
