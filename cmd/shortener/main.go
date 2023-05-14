package main

import (
	"log"
	"net/http"

	"github.com/darod1n/urlshorten/internal/config"
	"github.com/darod1n/urlshorten/internal/handlers"
	"github.com/darod1n/urlshorten/internal/storage"
	"github.com/go-chi/chi/v5"
)

func main() {
	serverConfig := config.NewConfig()

	router := chi.NewRouter()

	db := storage.NewDB()
	router.Get("/{shortURL}", func(w http.ResponseWriter, r *http.Request) {
		shortURL := chi.URLParam(r, "shortURL")
		handlers.GetBigURL(shortURL, db, w, r)
	})

	router.Post("/", func(w http.ResponseWriter, r *http.Request) {
		handlers.ShortURL(serverConfig.ServerHost, db, w, r)
	})

	http.Handle("/", router)

	err := http.ListenAndServe(serverConfig.Addr, nil)
	if err != nil {
		log.Fatal(err)
	}
}
