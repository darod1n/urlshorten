package main

import (
	"log"
	"net/http"

	"github.com/darod1n/urlshorten/internal/authorization"
	"github.com/darod1n/urlshorten/internal/compression"
	"github.com/darod1n/urlshorten/internal/config"
	"github.com/darod1n/urlshorten/internal/handlers"
	"github.com/darod1n/urlshorten/internal/logger"
	"github.com/darod1n/urlshorten/internal/storage"
	"github.com/go-chi/chi/v5"
)

const (
	HandlerIndex           = "/"
	HandlerShortURL        = "/{shortURL}"
	HandlerAPIShorten      = "/api/shorten"
	HandlerAPIShortenBatch = "/api/shorten/batch"
	HandlerAPIUserURLS     = "/api/user/urls"
	HandlerPing            = "/ping"
)

func main() {
	l, err := logger.InitializeLoger()
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}
	defer l.Sync()

	serverConfig := config.NewConfig()
	db, err := storage.NewDB(serverConfig)
	if err != nil {
		l.Fatalf("failed to create DB: %v", err)
	}
	defer db.Close()

	router := chi.NewRouter()
	router.Use(func(h http.Handler) http.Handler {
		return logger.WithLoggin(h, l)
	})

	router.Use(func(h http.Handler) http.Handler {
		return authorization.WithAutorization(h, serverConfig.SecretKey, l)
	})

	router.Use(compression.WithCompress)

	router.Get(HandlerShortURL, func(w http.ResponseWriter, r *http.Request) {
		shortURL := chi.URLParam(r, "shortURL")
		handlers.GetBigURL(shortURL, db, w, r)
	})
	router.Post(HandlerIndex, func(w http.ResponseWriter, r *http.Request) {
		handlers.ShortURL(serverConfig.ServerHost, db, w, r, l)
	})
	router.Post(HandlerAPIShorten, func(w http.ResponseWriter, r *http.Request) {
		handlers.APIShortenURL(serverConfig.ServerHost, db, w, r, l)
	})

	router.Post(HandlerAPIShortenBatch, func(w http.ResponseWriter, r *http.Request) {
		handlers.Batch(serverConfig.ServerHost, db, w, r, l)
	})

	router.Get(HandlerAPIUserURLS, func(w http.ResponseWriter, r *http.Request) {
		handlers.GetUserURLS(serverConfig.ServerHost, db, w, r, l)
	})

	router.Delete(HandlerAPIUserURLS, func(w http.ResponseWriter, r *http.Request) {
		handlers.DeleteUserURLS(db, w, r, l)
	})

	router.Get(HandlerPing, func(w http.ResponseWriter, r *http.Request) {
		handlers.Ping(db, w, r, l)
	})

	http.Handle(HandlerIndex, router)

	l.Infow(
		"Starting server",
		"Addr", serverConfig.Addr,
	)

	if err := http.ListenAndServe(serverConfig.Addr, nil); err != nil {
		l.Fatalf("failed to start server: %v", err)
	}
}
