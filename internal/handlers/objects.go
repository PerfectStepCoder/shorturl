package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PerfectStepCoder/shorturl/internal/models"
	"github.com/PerfectStepCoder/shorturl/internal/storage"
)

func ObjectShorterURL(mainStorage storage.Storage, baseURL string) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		// Декодирование запроса
		var requestFullURL models.RequestFullURL
		dec := json.NewDecoder(req.Body)
		if err := dec.Decode(&requestFullURL); err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		
		res.Header().Set("Content-Type", "application/json")

		shortURL, err := mainStorage.Save(requestFullURL.URL)
		if err != nil {
			var ue *storage.UniqURLError
			if errors.As(err, &ue) {
				originShortURL := strings.TrimSuffix(fmt.Sprintf("%s/%s", baseURL, ue.ShortHash), "\n")
				//log.Println("RESULT originShortURL:", originShortURL)
				res.Header().Set("Content-Type", "text/plain")
				res.WriteHeader(http.StatusConflict)
				fmt.Fprint(res, originShortURL)
				//http.Error(res, originShortURL, http.StatusConflict)
				return
			}
		}

		resp := models.ResponseShortURL{
			Result: strings.TrimSuffix(fmt.Sprintf("%s/%s", baseURL, shortURL), "\n"),
		}

		res.WriteHeader(http.StatusCreated)

		// Cериализуем ответ сервера
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Println("Error writing response:", err)
			return
		}
		//log.Println("RESULT:", jsonResp)

		res.Write(jsonResp)
	}
}

func ObjectsShorterURL(mainStorage storage.CorrelationStorage, baseURL string) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		// Декодирование запроса
		var requestCorrelationURLs []models.RequestCorrelationURL
		dec := json.NewDecoder(req.Body)
		if err := dec.Decode(&requestCorrelationURLs); err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		var correlationURLs []storage.CorrelationURL

		for _, value := range requestCorrelationURLs {
			correlationURLs = append(correlationURLs, storage.CorrelationURL{
				CorrelationID: value.CorrelationID,
				OriginalURL:   value.OriginalURL,
			})
		}

		shortURLs, err := mainStorage.CorrelationsSave(correlationURLs)

		if err != nil {
			var ue *storage.UniqURLError
			if errors.As(err, &ue) {
				originShortURL := strings.TrimSuffix(fmt.Sprintf("%s/%s", baseURL, ue.ShortHash), "\n")
				http.Error(res, originShortURL, http.StatusConflict)
				return
			}
		}

		// Кодирование ответа
		var resp []models.ResponseCorrelationURL
		for _, value := range shortURLs {
			resp = append(resp, models.ResponseCorrelationURL{
				CorrelationID: value, ShortURL: strings.TrimSuffix(fmt.Sprintf("%s/%s", baseURL, value), "\n"),
			})
		}

		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusCreated)

		// Cериализуем ответ сервера
		enc := json.NewEncoder(res)
		if err := enc.Encode(resp); err != nil {
			log.Println("Error writing response:", err)
			return
		}
	}
}
