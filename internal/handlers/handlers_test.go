package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockDB struct {
	urls map[string]string
}

func (db *MockDB) AddURL(url string, shortURL string) error {
	db.urls[shortURL] = url
	return nil
}

func (db *MockDB) GetURL(shortURL string) (string, bool) {
	bigURL, ok := db.urls[shortURL]
	return bigURL, ok
}

func TestShortURL(t *testing.T) {
	type want struct {
		code int
	}
	tests := []struct {
		name    string
		testURL string
		want    want
	}{
		{
			name: "postive test #1",
			want: want{
				code: 201,
			},
			testURL: `https://stackoverflow.com/questions/40096750/how-to-set-http-status-code-on-http-responsewriter`,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			bodyReader := strings.NewReader(test.testURL)

			request := httptest.NewRequest(http.MethodPost, "/", bodyReader)

			w := httptest.NewRecorder()
			db := &MockDB{urls: map[string]string{}}
			serverHost := "http://localhost:8080"
			ShortURL(serverHost, db, w, request, t)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, res.StatusCode, test.want.code)
		})
	}
}

func TestGetBigURL(t *testing.T) {
	type want struct {
		code     int
		location string
	}
	tests := []struct {
		name      string
		want      want
		testToken string
		addToken  bool
	}{
		{
			name: "positive test #1: existing short URL",
			want: want{
				code:     307,
				location: "https://stackoverflow.com/questions/40096750/how-to-set-http-status-code-on-http-responsewriter",
			},
			testToken: "hR3K3F",
			addToken:  true,
		},
		{
			name: "positive test #2: not existing short URL",
			want: want{
				code:     400,
				location: "",
			},
			testToken: "hR3K3F2",
			addToken:  false,
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/"+test.testToken, nil)

			w := httptest.NewRecorder()
			db := &MockDB{urls: map[string]string{}}

			if test.addToken {
				err := db.AddURL(test.want.location, test.testToken)
				if err != nil {
					t.Errorf("failed to add url: %v", err)
				}
			}

			GetBigURL(test.testToken, db, w, request)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, res.StatusCode, test.want.code)

		})
	}
}

func TestAPIShortenURL(t *testing.T) {
	type want struct {
		code int
	}
	tests := []struct {
		name    string
		want    want
		testURL string
	}{

		{
			name: "positive test #1",
			want: want{
				code: 201,
			},
			testURL: "https://stackoverflow.com/questions/40096750/how-to-set-http-status-code-on-http-responsewriter",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			jsonBody := map[string]interface{}{
				"url": test.testURL,
			}
			body, err := json.Marshal(jsonBody)

			if err != nil {
				t.Errorf("JSON Body error: %s", err)
			}

			request := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewReader(body))

			w := httptest.NewRecorder()
			db := &MockDB{urls: map[string]string{}}
			serverHost := "http://localhost:8080"
			APIShortenURL(serverHost, db, w, request, t)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, res.StatusCode, test.want.code)
		})
	}
}
