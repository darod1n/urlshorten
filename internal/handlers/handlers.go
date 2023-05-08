package handlers

import (
	"io"
	"net/http"

	"github.com/darod1n/urlshorten/internal/storage"

	"github.com/gorilla/mux"
)

type Storage interface {
	GenerateToken(lenToken int) string
	GetStrFromFile(path string) string
	SaveToFile(path string, str string)
	IsExist(token string) bool
}

func GenerateToken(lenToken int) string {
	return storage.GenerateToken(lenToken)
}

func GetStrFromFile(path string) string {
	return storage.GetStrFromFile(path)
}

func SaveToFile(path string, str string) {
	storage.SaveToFile(path, str)
}

func IsExist(token string) bool {
	return storage.IsExist(token)
}

func APIShortURL(serverHost string, db Storage, res http.ResponseWriter, req *http.Request) {
	body, _ := io.ReadAll(req.Body)
	token := db.GenerateToken(6)
	db.SaveToFile(token, string(body))
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(serverHost + req.RequestURI + token))
}

func APIGetBigURL(db Storage, res http.ResponseWriter, req *http.Request) {
	token := mux.Vars(req)["token"]
	if db.IsExist(token) {
		bigURL := db.GetStrFromFile(token)
		res.Header().Set("Location", bigURL)
		res.WriteHeader(http.StatusTemporaryRedirect)
	} else {
		res.WriteHeader(http.StatusBadRequest)
	}
}
