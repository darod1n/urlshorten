package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"strings"

	"github.com/darod1n/urlshorten/internal/helpers"
	"github.com/darod1n/urlshorten/internal/models"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type DB struct {
	base *sql.DB
}

const driverName = "pgx"

func (db *DB) AddURL(ctx context.Context, url string) (string, error) {
	shortURL := helpers.GenerateShortURL(10)
	row, err := db.base.ExecContext(ctx, "INSERT INTO urls (original_url, short_url) VALUES($1, $2) on conflict (original_url) do nothing;", url, shortURL)
	if err != nil {
		return "", err
	}

	rowAffected, err := row.RowsAffected()
	if err != nil {
		return "", err
	}

	if rowAffected == 0 {
		row := db.base.QueryRowContext(ctx, "SELECT short_url FROM urls where original_url=$1", url)
		var t string
		if err := row.Scan(&t); err != nil {
			return "", err
		}
		return t, models.ErrExistURL
	}

	return shortURL, nil
}

func (db *DB) GetURL(ctx context.Context, shortURL string) (string, error) {
	row := db.base.QueryRowContext(ctx, "select original_url from urls where short_url=$1;", shortURL)
	var originalURL string
	if err := row.Scan(&originalURL); err != nil {
		return "", err
	}
	return originalURL, nil
}

func (db *DB) PingContext(ctx context.Context) error {
	return db.base.PingContext(ctx)
}

func (db *DB) Close() error {
	return db.base.Close()
}

func (db *DB) Batch(ctx context.Context, host string, batch []models.BatchRequest) ([]models.BatchResponse, error) {

	var data []models.BatchResponse
	batchValues := make([]string, 0, len(batch))
	for _, val := range batch {
		shortURL := helpers.GenerateShortURL(10)
		valueQuery := fmt.Sprintf("('%s', '%s')", val.OriginURL, shortURL)
		batchValues = append(batchValues, valueQuery)

		url, err := url.JoinPath(host, shortURL)
		if err != nil {
			return nil, fmt.Errorf("failed to join path: %v", err)
		}
		data = append(data, models.BatchResponse{CorrelationID: val.CorrelationID, ShortURL: url})
	}
	query := fmt.Sprintf("INSERT INTO urls (original_url, short_url) VALUES %s on conflict (original_url) do nothing;", strings.Join(batchValues, ","))
	_, err := db.base.ExecContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to executed query: %v", err)
	}

	return data, nil
}

func createDB(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, "create table if not exists urls(short_url text primary key, original_url text unique);")
	if err != nil {
		return err
	}
	return nil
}

func NewDB(dataSourceName string) (*DB, error) {

	base, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	ctx := context.Background()
	if err := createDB(ctx, base); err != nil {
		return nil, fmt.Errorf("failed to create database: %v", err)
	}

	return &DB{
		base: base,
	}, nil
}
