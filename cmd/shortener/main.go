package main

import (
	"fmt"
	"net/http"

	"github.com/darod1n/urlshorten/internal/handlers"
	"github.com/darod1n/urlshorten/internal/storage"
	"github.com/gorilla/mux"
)

type Storage interface {
	GenerateToken(lenToken int) string
	GetStrFromFile(path string) string
	SaveToFile(path string, str string)
	IsExist(token string) bool
}

type DBStorage struct{}

func (db *DBStorage) GenerateToken(lenToken int) string {
	return storage.GenerateToken(lenToken)
}

func (db *DBStorage) GetStrFromFile(path string) string {
	return storage.GetStrFromFile(path)
}

func (db *DBStorage) SaveToFile(path string, str string) {
	storage.SaveToFile(path, str)
}

func (db *DBStorage) IsExist(token string) bool {
	return storage.IsExist(token)
}

func main() {
	router := mux.NewRouter()

	db := &DBStorage{}
	router.HandleFunc("/{token:[a-zA-z0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		handlers.APIGetBigURL(db, w, r)

	}).Methods("GET")
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handlers.APIShortURL(db, w, r)
	}).Methods("POST")

	http.Handle("/", router)

	err := http.ListenAndServe(`:8080`, nil)
	if err != nil {
		fmt.Println("Ошибочка")
		panic(err)
	}
}
