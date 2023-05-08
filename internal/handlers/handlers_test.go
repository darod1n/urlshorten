package handlers

import (
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

type storageMock struct {
	tokens map[string]string
}

func (db *storageMock) GenerateToken(lenToken int) string {
	var token string
	rand.Seed(time.Now().UnixNano())
	alphabet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	for i := 0; i < lenToken; i++ {
		r := rand.Intn(len(alphabet))
		token += string(alphabet[r])
	}
	return token
}

func (db *storageMock) GetStrFromFile(path string) string {
	return db.tokens[path]
}

func (db *storageMock) SaveToFile(path string, str string) {
	db.tokens[path] = str
}

func (db *storageMock) IsExist(token string) bool {
	if _, ok := db.tokens[token]; ok {
		return true
	} else {
		return false
	}
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
		t.Run(test.name, func(t *testing.T) {
			bodyReader := strings.NewReader(`https://stackoverflow.com/questions/40096750/how-to-set-http-status-code-on-http-responsewriter`)

			request := httptest.NewRequest(http.MethodPost, "/", bodyReader)

			w := httptest.NewRecorder()
			db := &storageMock{tokens: make(map[string]string)}

			APIShortURL(db, w, request)
			res := w.Result()

			assert.Equal(t, res.StatusCode, test.want.code)
			defer res.Body.Close()
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
			name: "positive test #1",
			want: want{
				code:     307,
				location: "https://stackoverflow.com/questions/40096750/how-to-set-http-status-code-on-http-responsewriter",
			},
			testToken: "/hR3K3F",
			addToken:  true,
		},
		{
			name: "positive test #2",
			want: want{
				code:     400,
				location: "",
			},
			testToken: "/hR3K3F2",
			addToken:  false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			request := httptest.NewRequest(http.MethodGet, test.testToken, nil)

			w := httptest.NewRecorder()
			db := &storageMock{tokens: make(map[string]string)}
			if test.addToken {
				db.SaveToFile(test.testToken, test.want.location)
			}

			tokens := map[string]string{
				"token": test.testToken,
			}
			request = mux.SetURLVars(request, tokens)

			APIGetBigURL(db, w, request)

			res := w.Result()
			assert.Equal(t, res.StatusCode, test.want.code)

			defer res.Body.Close()
		})
	}
}
