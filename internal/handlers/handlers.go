package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
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

type data struct {
	URL string `json:"url"`
}
type result struct {
	Result string `json:"result"`
}

func ShortURL(serverHost string, db Storage, res http.ResponseWriter, req *http.Request) {
	body, errBody := io.ReadAll(req.Body)
	if errBody != nil {
		log.Print(errBody)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	shortURL := helpers.GenerateShortURL(6)
	db.AddURL(string(body), shortURL)

	res.WriteHeader(http.StatusCreated)

	resultURL, errURL := url.JoinPath(serverHost, shortURL)

	if errURL != nil {
		log.Print(errURL)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	if _, errWrite := res.Write([]byte(resultURL)); errWrite != nil {
		log.Print(errWrite)
		res.WriteHeader(http.StatusBadRequest)
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

func APIShortenURL(serverHost string, db Storage, res http.ResponseWriter, req *http.Request) {
	var d data
	var result result

	var buf bytes.Buffer
	_, errBody := buf.ReadFrom(req.Body)
	if errBody != nil {
		return
	}

	if err := json.Unmarshal(buf.Bytes(), &d); err != nil {
		log.Print(err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Println(d, d.URL)
	shortURL := helpers.GenerateShortURL(6)
	fmt.Println(shortURL)
	db.AddURL(d.URL, shortURL)

	resultURL, errURL := url.JoinPath(serverHost, shortURL)
	if errURL != nil {
		log.Print(errURL)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	result.Result = resultURL
	ans, errJSON := json.Marshal(result)
	if errJSON != nil {
		log.Print(errJSON)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
	res.Write(ans)
}
