package handlers

import (
	"bytes"
	"encoding/json"
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

type Data struct {
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

/*
Запрос
POST http://localhost:8080/api/shorten HTTP/1.1
Host: localhost:8080
Content-Type: application/json

	{
	  "url": "https://practicum.yandex.ru"
	}

Ответ
HTTP/1.1 200 OK
Content-Type: application/json
Content-Length: 30

	{
	 "result": "http://localhost:8080/EwHXdJfB"
	}
*/
func APIShortenURL(db Storage, serverHost string, res http.ResponseWriter, req *http.Request) {
	var data Data

	var result result
	var buf bytes.Buffer
	_, errBody := buf.ReadFrom(req.Body)
	if errBody != nil {
		return
	}

	if err := json.Unmarshal(buf.Bytes(), &data); err != nil {
		return
	}

	bigURL := data.URL
	shortURL := helpers.GenerateShortURL(6)
	result.Result = serverHost + shortURL
	db.AddURL(bigURL, shortURL)

	ans, _ := json.Marshal(result)

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(ans))
}
