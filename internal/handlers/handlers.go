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
	if body, err := io.ReadAll(req.Body); err != nil {
		res.WriteHeader(http.StatusBadRequest)
	} else {
		shortURL := helpers.GenerateShortURL(6)
		db.AddURL(string(body), shortURL)
		res.WriteHeader(http.StatusCreated)
		resultURL, _ := url.JoinPath(serverHost, shortURL)
		_, err := res.Write([]byte(resultURL))
		if err != nil {
			log.Fatalf(err.Error())
		}
	}

}

func GetBigURL(shortURL string, db Storage, res http.ResponseWriter, req *http.Request) {
	if bigURL, ok := db.GetURL(shortURL); ok {
		res.Header().Set("Location", bigURL)
		res.WriteHeader(http.StatusTemporaryRedirect)
	} else {
		res.WriteHeader(http.StatusBadRequest)
	}
}
