package handlers

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"github.com/PerfectStepCoder/shorturl/internal/storage"
	"github.com/go-chi/chi/v5"
)

func ShorterURL(mainStorage storage.Storage, baseURL string) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		originURLbytes, _ := io.ReadAll(req.Body)
		defer func() {
			if err := req.Body.Close(); err != nil {
				log.Printf("could not close response body: %v", err)
			}
		}()
		originURL := string(originURLbytes)
		if originURL == "" {
			http.Error(res, "URL not send", http.StatusBadRequest)
			return
		}
		shortURL, err := mainStorage.Save(originURL)
		if err != nil {
			var ue *storage.UniqURLError
			if errors.As(err, &ue) {
				originShortURL := strings.TrimSuffix(fmt.Sprintf("%s/%s", baseURL, ue.ShortHash), "\n")
				http.Error(res, originShortURL, http.StatusConflict)
				return
			}
		}
		shortURLfull := fmt.Sprintf("%s/%s", baseURL, shortURL)
		res.WriteHeader(http.StatusCreated)
		res.Header().Set("content-type", "text/plain")
		if _, err := res.Write([]byte(shortURLfull)); err != nil {
			log.Println("Error writing response:", err)
		}
	}
}

func GetURL(storage storage.Storage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
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
