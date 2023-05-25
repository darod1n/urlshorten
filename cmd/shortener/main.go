package main

import (
	"log"
	"net/http"

	"github.com/darod1n/urlshorten/internal/compression"
	"github.com/darod1n/urlshorten/internal/config"
	"github.com/darod1n/urlshorten/internal/handlers"
	"github.com/darod1n/urlshorten/internal/logger"
	"github.com/darod1n/urlshorten/internal/storage"
	"github.com/go-chi/chi/v5"
)

func main() {

	serverConfig := config.NewConfig()
	db := storage.NewDB()
	if err := logger.InitializeLoger(); err != nil {
		log.Fatal(err)
	}

	router := chi.NewRouter()
	router.Use(logger.WithLoggin)
	router.Use(compression.WithCompress)
	router.Get("/{shortURL}", func(w http.ResponseWriter, r *http.Request) {
		shortURL := chi.URLParam(r, "shortURL")
		handlers.GetBigURL(shortURL, db, w, r)
	})
	router.Post("/", func(w http.ResponseWriter, r *http.Request) {
		handlers.ShortURL(serverConfig.ServerHost, db, w, r)
	})
	router.Post("/api/shorten", func(w http.ResponseWriter, r *http.Request) {
		handlers.APIShortenURL(serverConfig.ServerHost, db, w, r)
	})

	http.Handle("/", router)

	logger.Sugar.Infow(
		"Starting server",
		"Addr", serverConfig.Addr,
	)

	errServer := http.ListenAndServe(serverConfig.Addr, nil)
	if errServer != nil {
		log.Fatal(errServer)
	}
}
