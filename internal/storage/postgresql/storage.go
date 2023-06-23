package postgresql

import (
	"context"
	"fmt"
	"net/url"

	"github.com/darod1n/urlshorten/internal/helpers"
	"github.com/darod1n/urlshorten/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	base *pgxpool.Pool
}

const driverName = "pgx"

func (db *DB) AddURL(ctx context.Context, url string) (string, error) {
	shortURL := helpers.GenerateShortURL(url, 10)
	row := db.base.QueryRow(ctx, `
	with dataNew as (
		select 
			$2 as short_url,
			$1 as original_url
	), 
	dupData as (
	select 
			short_url,
			original_url
		from urls
		where original_url  in (select original_url  from dataNew)
	), 
	
	insData as (
		insert into urls
			select * from dataNew
			where (original_url) not in (select original_url from dupData)
			returning short_url
	---
	) 
	
	select short_url, true as result from dupData union all select short_url, false as result from insData;
	`, url, shortURL)
	var queryShortURL string
	var status bool
	if err := row.Scan(&queryShortURL, &status); err != nil {
		return "", fmt.Errorf("failed scan query: %v", err)
	}
	if status {
		return queryShortURL, models.ErrExistURL
	}

	return shortURL, nil
}

func (db *DB) GetURL(ctx context.Context, shortURL string) (string, error) {
	row := db.base.QueryRow(ctx, "select original_url from urls where short_url=$1;", shortURL)
	var originalURL string
	if err := row.Scan(&originalURL); err != nil {
		return "", fmt.Errorf("failed to scan query row: %v", err)
	}
	return originalURL, nil
}

func (db *DB) PingContext(ctx context.Context) error {
	return db.base.Ping(ctx)
}

func (db *DB) Close() {
	db.base.Close()
}

func (db *DB) Batch(ctx context.Context, host string, br []models.BatchRequest) ([]models.BatchResponse, error) {

	batch := &pgx.Batch{}
	tx, err := db.base.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin tx: %v", err)
	}
	defer tx.Commit(ctx)

	var data []models.BatchResponse
	for _, val := range br {

		shortURL := helpers.GenerateShortURL(val.OriginURL, 10)
		batch.Queue("INSERT INTO urls (original_url, short_url) VALUES ($1, $2) on conflict (original_url) do nothing;", val.OriginURL, shortURL)
		url, err := url.JoinPath(host, shortURL)
		if err != nil {
			return nil, fmt.Errorf("failed to join path: %v", err)
		}
		data = append(data, models.BatchResponse{CorrelationID: val.CorrelationID, ShortURL: url})
	}
	b := tx.SendBatch(ctx, batch)
	defer b.Close()

	if _, err := b.Exec(); err != nil {
		return nil, fmt.Errorf("failed to executed query: %v", err)
	}

	return data, nil
}

func createDB(ctx context.Context, db *pgxpool.Pool) error {
	_, err := db.Exec(ctx, "create table if not exists urls(short_url text primary key, original_url text unique);")
	if err != nil {
		return err
	}
	return nil
}

func NewDB(dataSourceName string) (*DB, error) {
	ctx := context.Background()
	base, err := pgxpool.New(ctx, dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	if err := createDB(ctx, base); err != nil {
		return nil, fmt.Errorf("failed to create database: %v", err)
	}

	return &DB{
		base: base,
	}, nil
}
