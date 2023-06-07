package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/darod1n/urlshorten/internal/helpers"
)

type Storage interface {
	AddURL(url string, shortURL string) error
	GetURL(shortURL string) (string, bool)
}

type logger interface {
	Errorf(template string, args ...interface{})
}

type data struct {
	URL string `json:"url"`
}
type result struct {
	Result string `json:"result"`
}

func ShortURL(serverHost string, db Storage, res http.ResponseWriter, req *http.Request, l logger) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		l.Errorf("failed to read body: %v", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	shortURL := helpers.GenerateShortURL(6)
	errAddURL := db.AddURL(string(body), shortURL)
	if errAddURL != nil {
		l.Errorf("failed to add url: %v", errAddURL)
		res.WriteHeader((http.StatusBadRequest))
		return
	}

	res.WriteHeader(http.StatusCreated)

	resultURL, errURL := url.JoinPath(serverHost, shortURL)

	if errURL != nil {
		l.Errorf("failed to join path: %v", errURL)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	if _, errWrite := res.Write([]byte(resultURL)); errWrite != nil {
		l.Errorf("failed to write byte: %v", errURL)
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

func APIShortenURL(serverHost string, db Storage, res http.ResponseWriter, req *http.Request, l logger) {
	var buf bytes.Buffer
	_, errBody := buf.ReadFrom(req.Body)
	if errBody != nil {
		l.Errorf("failed to read body: %v", errBody)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	var d data
	if err := json.Unmarshal(buf.Bytes(), &d); err != nil {
		l.Errorf("failed to unmarshal the request body: %v", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	shortURL := helpers.GenerateShortURL(6)
	db.AddURL(d.URL, shortURL)

	resultURL, errURL := url.JoinPath(serverHost, shortURL)
	if errURL != nil {
		l.Errorf("failed to join path: %v", errURL)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	var result result
	result.Result = resultURL
	ans, errJSON := json.Marshal(result)
	if errJSON != nil {
		l.Errorf("failed to marshal result: %v", errJSON)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
	res.Write(ans)
}
