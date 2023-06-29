package memory

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"sync"

	"github.com/darod1n/urlshorten/internal/helpers"
	"github.com/darod1n/urlshorten/internal/models"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type DB struct {
	urls map[string]string
	mu   *sync.Mutex
}

func (db *DB) AddURL(ctx context.Context, url string) (string, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	shortURL := helpers.GenerateShortURL(url, 10)
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
	return models.ErrNYI
}

func (db *DB) Close() {
}

func (db *DB) GetUserURLS(ctx context.Context, host string) ([]models.UserURLS, error) {
	return nil, fmt.Errorf("not yet implemented")
}

func (db *DB) Batch(ctx context.Context, host string, batch []models.BatchRequest) ([]models.BatchResponse, error) {
	br := make([]models.BatchResponse, len(batch))
	for _, val := range batch {
		shortURL, err := db.AddURL(ctx, val.OriginURL)
		if err != nil {
			return nil, fmt.Errorf("failed to add url: %w", err)
		}
		url, err := url.JoinPath(host, shortURL)
		if err != nil {
			return nil, fmt.Errorf("failed to join path: %w", err)
		}

		br = append(br, models.BatchResponse{CorrelationID: val.CorrelationID, ShortURL: url})
	}
	return br, nil
}
func (db *DB) DeleteUserURLS(ctx context.Context, userID string, urls []string) error {
	return models.ErrRemoveURL
}

func NewDB(urls map[string]string) (*DB, error) {
	return &DB{
		urls: urls,
		mu:   &sync.Mutex{},
	}, nil
}
