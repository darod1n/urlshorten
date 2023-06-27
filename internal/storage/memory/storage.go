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

func (db *DB) GetURL(ctx context.Context, shortURL string) (string, bool, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	bigURL, ok := db.urls[shortURL]
	if ok {
		return bigURL, false, nil
	}
	return "", false, errors.New("failed to get url")
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
			return nil, err
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

func NewDB(urls map[string]string) (*DB, error) {
	return &DB{
		urls: urls,
		mu:   &sync.Mutex{},
	}, nil
}
