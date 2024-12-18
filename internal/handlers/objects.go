// Модуль содержит декораторы для обработки запросов авторизованных пользователей.
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

// ObjectShorterURL - обработка одной ссылоки за один запрос.
func ObjectShorterURL(mainStorage storage.Storage, baseURL string) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		// Аутентификация
		var userUID string
		cookies, err := req.Cookie("userUID")
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				log.Print("Cookie 'userUID' отсутствует, создается новый.")
				userUID, _ = SetNewCookie(res)
			} else {
				// Обработка других возможных ошибок
				log.Printf("Ошибка при получении cookie: %v", err)
				http.Error(res, "Ошибка сервера", http.StatusInternalServerError)
				return
			}
		} else {
			// Если cookie существует, выполняем валидацию
			userUID, _ = ValidateUserUID(cookies.Value) // обработка исключения не требуется
		}

		// Декодирование запроса
		var requestFullURL models.RequestFullURL
		dec := json.NewDecoder(req.Body)
		if err := dec.Decode(&requestFullURL); err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		res.Header().Set("Content-Type", "application/json")

		shortURL, err := mainStorage.Save(requestFullURL.URL, userUID)
		if err != nil {
			var ue *storage.UniqURLError
			if errors.As(err, &ue) {
				originShortURL := strings.TrimSuffix(fmt.Sprintf("%s/%s", baseURL, ue.ShortHash), "\n")
				res.WriteHeader(http.StatusConflict)
				resp := models.ResponseShortURL{
					Result: originShortURL,
				}
				// Cериализуем ответ сервера
				jsonResp, err := json.Marshal(resp)
				if err != nil {
					log.Println("Error writing response:", err)
					return
				}
				res.Write(jsonResp)
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

		res.Write(jsonResp)
	}
}

// ObjectsShorterURL - обработка несколько ссылок за один запрос.
func ObjectsShorterURL(mainStorage storage.CorrelationStorage, baseURL string) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		// Аутентификация
		var userUID string
		cookies, err := req.Cookie("userUID")
		if err != nil {
			log.Print("No cookies")
			userUID, _ = SetNewCookie(res)
		} else {
			userUID, _ = ValidateUserUID(cookies.Value) // обработка исключения не требуется
		}

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

		shortURLs, err := mainStorage.CorrelationsSave(correlationURLs, userUID)

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
