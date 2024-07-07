package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/PerfectStepCoder/shorturl/internal/storage"
)


func ShorterURL(storage *storage.Storage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		originURLbytes, _ := io.ReadAll(req.Body)
		originURL := string(originURLbytes)
		shortURL := storage.Save(originURL)
		shortURLfull := fmt.Sprintf("%s/%s", "http://localhost:8080", shortURL) 
		res.WriteHeader(http.StatusCreated)
		res.Header().Set("content-type", "text/plain")
		if _, err := res.Write([]byte(shortURLfull)); err != nil {
			log.Println("Error writing response:", err)
		}
	}
}

func GetURL(storage *storage.Storage) http.HandlerFunc {
	return func (res http.ResponseWriter, req *http.Request) {
		shortURL := req.PathValue("id")
		originURL, exist := storage.Get(shortURL)
		if !exist {
			http.Error(res, "Not Found", http.StatusNotFound)
		}
		res.Header().Set("Location", originURL)
		res.WriteHeader(http.StatusTemporaryRedirect)
	}
}

