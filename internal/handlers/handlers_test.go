package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockDB struct {
	urls map[string]string
}

func (db *MockDB) AddURL(url string, shortURL string) {
	db.urls[shortURL] = url
}

func (db *MockDB) GetURL(shortURL string) (string, bool) {
	bigURL, ok := db.urls[shortURL]
	return bigURL, ok
}

func TestAPIShortURL(t *testing.T) {
	type want struct {
		code int
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "postive test #1",
			want: want{
				code: 201,
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			bodyReader := strings.NewReader(`https://stackoverflow.com/questions/40096750/how-to-set-http-status-code-on-http-responsewriter`)

			request := httptest.NewRequest(http.MethodPost, "/", bodyReader)

			w := httptest.NewRecorder()
			db := &MockDB{urls: map[string]string{}}
			serverHost := "http://localhost:8080"
			APIShortURL(serverHost, db, w, request)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, res.StatusCode, test.want.code)

		})
	}
}

func TestAPIGetBigURL(t *testing.T) {
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
				db.AddURL(test.want.location, test.testToken)
			}

			APIGetBigURL(test.testToken, db, w, request)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, res.StatusCode, test.want.code)

		})
	}
}
