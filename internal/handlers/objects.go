package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/PerfectStepCoder/shorturl/internal/models"
	"github.com/PerfectStepCoder/shorturl/internal/storage"
)

func ObjectShorterURL(storage storage.Storage, baseURL string) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		// Декодирование запроса
		var requestFullURL models.RequestFullURL
		dec := json.NewDecoder(req.Body)
		if err := dec.Decode(&requestFullURL); err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		shortURL := storage.Save(requestFullURL.URL)
		resp := models.ResponseShortURL{
			Result: fmt.Sprintf("%s/%s", baseURL, shortURL),
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
