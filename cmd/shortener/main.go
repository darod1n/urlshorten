package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/darod1n/urlshorten/internal/compression"
	"github.com/darod1n/urlshorten/internal/config"
	"github.com/darod1n/urlshorten/internal/handlers"
	"github.com/darod1n/urlshorten/internal/logger"
	"github.com/darod1n/urlshorten/internal/storage"
	"github.com/go-chi/chi/v5"
)

func main() {
	l, err := logger.InitializeLoger()
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}
	defer l.Sync()

	serverConfig := config.NewConfig()
	db, err := storage.NewDB(serverConfig.Path, serverConfig.DriverName, serverConfig.DataSourceName)
	if err != nil {
		l.Fatalf("failed to create DB: %v", err)
	}

	router := chi.NewRouter()
	router.Use(func(h http.Handler) http.Handler {
		return logger.WithLoggin(h, l)
	})
	router.Use(compression.WithCompress)
	router.Get("/{shortURL}", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		shortURL := chi.URLParam(r, "shortURL")
		handlers.GetBigURL(ctx, shortURL, db, w, r)
	})
	router.Post("/", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		handlers.ShortURL(ctx, serverConfig.ServerHost, db, w, r, l)
	})
	router.Post("/api/shorten", func(w http.ResponseWriter, r *http.Request) {
		handlers.APIShortenURL(serverConfig.ServerHost, db, w, r, l)
	})
	router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		handlers.Ping(ctx, db, w, r, l)
	})

	http.Handle("/", router)

	l.Infow(
		"Starting server",
		"Addr", serverConfig.Addr,
	)

	if err := http.ListenAndServe(serverConfig.Addr, nil); err != nil {
		l.Fatalf("failed to start server: %v", err)
	}
}
