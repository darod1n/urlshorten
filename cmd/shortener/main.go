package main

import (
	"fmt"
	"net/http"

	"github.com/darod1n/urlshorten/internal/handlers"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/{token:[a-zA-z0-9]+}", handlers.APIGetBigURL).Methods("GET")
	router.HandleFunc("/", handlers.APIShortURL).Methods("POST")

	http.Handle("/", router)

	err := http.ListenAndServe(`:8080`, nil)
	if err != nil {
		fmt.Println("Ошибочка")
		panic(err)
	}
}
