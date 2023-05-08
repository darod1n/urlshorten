package handlers

import (
	"io"
	"net/http"
	"os"

	"github.com/darod1n/urlshorten/internal/util"

	"github.com/gorilla/mux"
)

func APIShortURL(res http.ResponseWriter, req *http.Request) {
	body, _ := io.ReadAll(req.Body)
	token := util.GenerateToken(6)
	util.SaveToFile(token, string(body))
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte("http://localhost:8080/" + token))
}

func APIGetBigURL(res http.ResponseWriter, req *http.Request) {
	token := mux.Vars(req)["token"]
	dir, _ := os.Getwd()
	if _, err := os.Stat(dir + "/secrets/" + token); os.IsNotExist(err) {
		res.WriteHeader(400)
	} else {
		bigURL := util.GetStrFromFile(dir + "/secrets/" + token)
		res.Header().Set("Location", bigURL)
		res.WriteHeader(307)
	}

}
