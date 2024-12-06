// Модуль содержит декораторы для обработки запросов авторизованных пользователей.
package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/PerfectStepCoder/shorturl/internal/models"
	"github.com/PerfectStepCoder/shorturl/internal/storage"
)

// ShorterStats - обрабатывает запрос к /api/internal/stats
func ShorterStats(mainStorage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Формирование и отправка ответа
		w.Header().Set("Content-Type", "application/json")

		users, _ := mainStorage.CountUsers()
		urls, _ := mainStorage.CountURLs()
		resp := models.ResponseStatsBase{
			Urls:  urls,
			Users: users,
		}

		// Cериализуем ответ сервера
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Println("Error writing response:", err)
			return
		}

		w.Write(jsonResp)
		w.WriteHeader(http.StatusOK)
	}
}
