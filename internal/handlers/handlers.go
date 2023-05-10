package handlers

import (
	"io"
	"net/http"

	"github.com/darod1n/urlshorten/internal/helpers"
	"github.com/darod1n/urlshorten/internal/storage"
)

type DB struct {
	storage.DB
}

type Storage interface {
	storage.Storage
}

func APIShortURL(serverHost string, db Storage, res http.ResponseWriter, req *http.Request) {
	if body, err := io.ReadAll(req.Body); err != nil {
		res.WriteHeader(http.StatusBadRequest)
	} else {
		shortURL := helpers.GenerateShortURL(6)
		db.AddURL(string(body), shortURL)
		res.WriteHeader(http.StatusCreated)
		res.Write([]byte(serverHost + req.RequestURI + shortURL))
	}

}

func APIGetBigURL(db Storage, res http.ResponseWriter, req *http.Request) {
	shortURL := req.RequestURI[1:]
	if bigURL, ok := db.GetURL(shortURL); ok {
		res.Header().Set("Location", bigURL)
		res.WriteHeader(http.StatusTemporaryRedirect)
	} else {
		res.WriteHeader(http.StatusBadRequest)
	}
}
