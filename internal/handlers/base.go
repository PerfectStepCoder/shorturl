package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"github.com/go-chi/chi/v5"
	"github.com/PerfectStepCoder/shorturl/internal/storage"
)


func ShorterURL(storage *storage.Storage, baseURL string) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		originURLbytes, _ := io.ReadAll(req.Body)
		originURL := string(originURLbytes)
		if originURL == "" {
			http.Error(res, "URL not send", http.StatusBadRequest)
			return
		}
		shortURL := storage.Save(originURL)
		shortURLfull := fmt.Sprintf("%s/%s", baseURL, shortURL) 
		res.WriteHeader(http.StatusCreated)
		res.Header().Set("content-type", "text/plain")
		if _, err := res.Write([]byte(shortURLfull)); err != nil {
			log.Println("Error writing response:", err)
		}
	}
}

func GetURL(storage *storage.Storage) http.HandlerFunc {
	return func (res http.ResponseWriter, req *http.Request) {
		shortURL := chi.URLParam(req, "id")
		if shortURL == "" {
			http.Error(res, "ShortURL not send", http.StatusBadRequest)
			return
		}
		originURL, exists := storage.Get(shortURL)
		if !exists {
			http.Error(res, "Not Found", http.StatusNotFound)
			return
		}
		res.Header().Set("Location", originURL)
		res.WriteHeader(http.StatusTemporaryRedirect)
	}
}

