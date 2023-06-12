package storage

import (
	"context"

	"github.com/darod1n/urlshorten/internal/storage/file"
	"github.com/darod1n/urlshorten/internal/storage/memory"
	"github.com/darod1n/urlshorten/internal/storage/postgresql"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type DB interface {
	AddURL(ctx context.Context, url string) (string, error)
	GetURL(ctx context.Context, shortURL string) (string, error)
	PingContext(ctx context.Context) error
}

func NewDB(path, driverName, dataSourceName string) (DB, error) {
	if dataSourceName != "" {
		return postgresql.NewDB(path, driverName, dataSourceName)
	}

	if path != "" {
		return file.NewDB(path)
	}
	return memory.NewDB()
}
