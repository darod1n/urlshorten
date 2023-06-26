package storage

import (
	"context"

	"github.com/darod1n/urlshorten/internal/config"
	"github.com/darod1n/urlshorten/internal/models"
	"github.com/darod1n/urlshorten/internal/storage/file"
	"github.com/darod1n/urlshorten/internal/storage/memory"
	"github.com/darod1n/urlshorten/internal/storage/postgresql"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type DB interface {
	AddURL(ctx context.Context, url string) (string, error)
	GetURL(ctx context.Context, shortURL string) (string, error)
	PingContext(ctx context.Context) error
	Batch(ctx context.Context, host string, batch []models.BatchRequest) ([]models.BatchResponse, error)
	CreateUserID(ctx context.Context) (string, error)
	GetUserURLS(ctx context.Context) ([]models.UserURLS, error)
	Close()
}

func NewDB(cfg *config.Config) (DB, error) {
	if cfg.DataSourceName != "" {
		return postgresql.NewDB(cfg.DataSourceName)
	}

	if cfg.Path != "" {
		return file.NewDB(cfg.Path)
	}
	urls := make(map[string]string)
	return memory.NewDB(urls)
}
