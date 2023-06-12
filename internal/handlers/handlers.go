package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

type Storage interface {
	AddURL(ctx context.Context, url string) (string, error)
	GetURL(ctx context.Context, shortURL string) (string, error)
	PingContext(ctx context.Context) error
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

func ShortURL(ctx context.Context, serverHost string, db Storage, res http.ResponseWriter, req *http.Request, l logger) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		l.Errorf("failed to read body: %v", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	shortURL, err := db.AddURL(ctx, string(body))

	if err != nil {
		l.Errorf("failed to add url: %v", err)
		res.WriteHeader((http.StatusBadRequest))
		return
	}

	res.WriteHeader(http.StatusCreated)

	resultURL, err := url.JoinPath(serverHost, shortURL)

	if err != nil {
		l.Errorf("failed to join path: %v", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	if _, err := res.Write([]byte(resultURL)); err != nil {
		l.Errorf("failed to write byte: %v", err)
		res.WriteHeader(http.StatusBadRequest)
	}
}

func GetBigURL(ctx context.Context, shortURL string, db Storage, res http.ResponseWriter, req *http.Request) {
	bigURL, err := db.GetURL(ctx, shortURL)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	res.Header().Set("Location", bigURL)
	res.WriteHeader(http.StatusTemporaryRedirect)
}

func APIShortenURL(serverHost string, db Storage, res http.ResponseWriter, req *http.Request, l logger) {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		l.Errorf("failed to read body: %v", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	var d data
	if err := json.Unmarshal(buf.Bytes(), &d); err != nil {
		l.Errorf("failed to unmarshal the request body: %v", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx := context.Background()

	shortURL, err := db.AddURL(ctx, d.URL)
	if err != nil {
		l.Errorf("failed to add url: %v", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	resultURL, err := url.JoinPath(serverHost, shortURL)
	if err != nil {
		l.Errorf("failed to join path: %v", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	var result result
	result.Result = resultURL
	ans, err := json.Marshal(result)
	if err != nil {
		l.Errorf("failed to marshal result: %v", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
	res.Write(ans)
}

func Ping(ctx context.Context, db Storage, res http.ResponseWriter, req *http.Request, l logger) {
	if err := db.PingContext(ctx); err != nil {
		l.Errorf("failed to ping database: %v", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusOK)
}
