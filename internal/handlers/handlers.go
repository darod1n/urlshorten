package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/darod1n/urlshorten/internal/models"
	"github.com/darod1n/urlshorten/internal/storage/errstorage"
)

type Storage interface {
	AddURL(ctx context.Context, url string) (string, error)
	GetURL(ctx context.Context, shortURL string) (string, error)
	PingContext(ctx context.Context) error
	Batch(ctx context.Context, host string, batch []models.BatchRequest) ([]models.BatchResponse, error)
	GetUserURLS(ctx context.Context, host string) ([]models.UserURLS, error)
	DeleteUserURLS(ctx context.Context, userID string, urls []string) error
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

	status := http.StatusCreated
	ctx := req.Context()
	shortURL, err := db.AddURL(ctx, string(body))
	if err != nil {
		if !errors.Is(err, errstorage.ExistURL) {
			l.Errorf("failed to add url: %v", err)
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		status = http.StatusConflict
	}

	resultURL, err := url.JoinPath(serverHost, shortURL)
	if err != nil {
		l.Errorf("failed to join path: %v", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	res.WriteHeader(status)
	if _, err := res.Write([]byte(resultURL)); err != nil {
		l.Errorf("failed to write byte: %v", err)
		res.WriteHeader(http.StatusBadRequest)
	}
}

func GetBigURL(shortURL string, db Storage, res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	bigURL, err := db.GetURL(ctx, shortURL)
	if err != nil {
		if errors.Is(err, errstorage.RemoveURL) {
			res.WriteHeader(http.StatusGone)
			return
		}
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

	ctx := req.Context()
	status := http.StatusCreated
	var result result
	shortURL, err := db.AddURL(ctx, d.URL)
	if err != nil {
		if shortURL == "" {
			l.Errorf("failed to add url: %v", err)
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		status = http.StatusConflict
	}

	resultURL, err := url.JoinPath(serverHost, shortURL)
	if err != nil {
		l.Errorf("failed to join path: %v", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	result.Result = resultURL
	ans, err := json.Marshal(result)
	if err != nil {
		l.Errorf("failed to marshal result: %v", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(status)
	res.Write(ans)
}

func Batch(serverHost string, db Storage, res http.ResponseWriter, req *http.Request, l logger) {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		l.Errorf("failed to read body: %v", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	var bodyBatch []models.BatchRequest
	if err := json.Unmarshal(buf.Bytes(), &bodyBatch); err != nil {
		l.Errorf("failed to unmarshal body: %v", err)
		res.WriteHeader(http.StatusBadRequest)
		return

	}

	ctx := req.Context()
	resp, err := db.Batch(ctx, serverHost, bodyBatch)
	if err != nil {
		l.Errorf("failed to batch: %v", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	ans, err := json.Marshal(resp)
	if err != nil {
		l.Errorf("failed to marshal result: %v", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
	res.Write(ans)
}

func Ping(db Storage, res http.ResponseWriter, req *http.Request, l logger) {
	ctx := req.Context()
	if err := db.PingContext(ctx); err != nil {
		l.Errorf("failed to ping database: %v", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusOK)
}

func GetUserURLS(serverHost string, db Storage, res http.ResponseWriter, req *http.Request, l logger) {
	ctx := req.Context()
	userURLS, err := db.GetUserURLS(ctx, serverHost)
	if err != nil {
		l.Errorf("failed to get user urls: %v", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	if len(userURLS) == 0 {
		res.WriteHeader(http.StatusNoContent)
		return
	}

	ans, err := json.Marshal(userURLS)
	if err != nil {
		l.Errorf("failed to marshal result: %v", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(ans)
}

func DeleteUserURLS(db Storage, res http.ResponseWriter, req *http.Request, l logger) {
	var buf bytes.Buffer

	if _, err := buf.ReadFrom(req.Body); err != nil {
		l.Errorf("failed to read body: %v", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	var shortURLS []string
	if err := json.Unmarshal(buf.Bytes(), &shortURLS); err != nil {
		l.Errorf("failed to unmarsharl: %v", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	userID := req.Context().Value(models.CtxKeyUserID)

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		if err := db.DeleteUserURLS(ctx, userID.(string), shortURLS); err != nil {
			l.Errorf("failed delete user urls: %v", err)
		}
	}()

	res.WriteHeader(http.StatusAccepted)
}
