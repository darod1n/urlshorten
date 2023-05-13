package main

import (
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
	router.Get("/{shortURL:[a-zA-Z0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		shortURL := chi.URLParam(r, "shortURL")
		handlers.APIGetBigURL(shortURL, db, w, r)
	})

	router.Post("/", func(w http.ResponseWriter, r *http.Request) {
		handlers.APIShortURL(serverConfig.ServerHost, db, w, r)
	})

	http.Handle("/", router)

	err := http.ListenAndServe(serverConfig.Addr, nil)
	if err != nil {
		panic(err)
	}
}
