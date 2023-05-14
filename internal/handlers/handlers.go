package handlers

import (
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/darod1n/urlshorten/internal/helpers"
)

type Storage interface {
	AddURL(url string, shortURL string)
	GetURL(shortURL string) (string, bool)
}

func ShortURL(serverHost string, db Storage, res http.ResponseWriter, req *http.Request) {
	body, errBody := io.ReadAll(req.Body)
	if errBody != nil {
		log.Fatal(errBody)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	shortURL := helpers.GenerateShortURL(6)
	db.AddURL(string(body), shortURL)

	res.WriteHeader(http.StatusCreated)

	resultURL, errURL := url.JoinPath(serverHost, shortURL)

	if errURL != nil {
		log.Fatal(errURL)
		return
	}

	if _, errWrite := res.Write([]byte(resultURL)); errWrite != nil {
		log.Fatal(errWrite)
	}
}

func GetBigURL(shortURL string, db Storage, res http.ResponseWriter, req *http.Request) {
	bigURL, ok := db.GetURL(shortURL)
	if !ok {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	res.Header().Set("Location", bigURL)
	res.WriteHeader(http.StatusTemporaryRedirect)
}
