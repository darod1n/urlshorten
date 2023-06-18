package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"net/url"

	"github.com/darod1n/urlshorten/internal/models"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type DB struct {
	base *sql.DB
}

const driverName = "pgx"

func (db *DB) AddURL(ctx context.Context, url string) (string, error) {
	row := db.base.QueryRowContext(ctx, "INSERT INTO urls (original_url) VALUES($1) on conflict (original_url) do nothing returning short_url;", url)

	var shortURL string
	row.Scan(&shortURL)

	if shortURL == "" {
		row := db.base.QueryRowContext(ctx, "select short_url from urls where original_url=$1;", url)
		row.Scan(&shortURL)
		return shortURL, errors.New("origin url is exist")
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
	tx, err := db.base.Begin()
	if err != nil {
		return nil, err
	}

	var data []models.BatchResponse
	for _, val := range batch {
		row := tx.QueryRowContext(ctx, "INSERT INTO urls (original_url) VALUES($1) returning short_url;", val.OriginURL)

		var shortURL string
		err := row.Scan(&shortURL)
		if err != nil {
			return nil, err
		}

		url, _ := url.JoinPath(host, shortURL)
		data = append(data, models.BatchResponse{CorrelationID: val.CorrelationID, ShortURL: url})
	}
	return data, tx.Commit()
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
		return nil, err
	}

	ctx := context.Background()
	if err := createDB(ctx, base); err != nil {
		return nil, err
	}

	return &DB{
		base: base,
	}, nil
}
