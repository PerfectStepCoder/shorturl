// Модуль handlers содержит обработчики запросов HTTP сервиса.
package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/PerfectStepCoder/shorturl/internal/models"
	"github.com/PerfectStepCoder/shorturl/internal/storage"
	"github.com/go-chi/chi/v5"
)

// batchSize - размер батча для массового удаления ссылок.
const batchSize = 15

// ShorterURL - обработчик ссылок.
func ShorterURL(mainStorage storage.Storage, baseURL string) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		// Аутентификация
		userUID := fmt.Sprintf("%s", req.Context().Value(UserKeyUID))

		originURLbytes, _ := io.ReadAll(req.Body)

		originURL := string(originURLbytes)
		if originURL == "" {
			http.Error(res, "URL not send", http.StatusBadRequest)
			return
		}
		shortURL, err := mainStorage.Save(originURL, userUID)
		if err != nil {
			var ue *storage.UniqURLError
			if errors.As(err, &ue) {
				originShortURL := strings.TrimSuffix(fmt.Sprintf("%s/%s", baseURL, ue.ShortHash), "\n")
				res.WriteHeader(http.StatusConflict)
				res.Write([]byte(originShortURL))
				return
			}
		}
		shortURLfull := strings.TrimSuffix(fmt.Sprintf("%s/%s", baseURL, shortURL), "\n")
		res.WriteHeader(http.StatusCreated)
		res.Header().Set("Content-Type", "application/json")
		if _, err := res.Write([]byte(shortURLfull)); err != nil {
			log.Println("Error writing response:", err)
		}
	}
}

// GetURL - возвращает оригинальную ссылку по передаваемой сокращенной ссылке.
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
		result, _ := storage.IsDeleted(shortURL)
		if result {
			res.WriteHeader(http.StatusGone)
			return
		}
		res.Header().Set("Location", originURL)
		res.WriteHeader(http.StatusTemporaryRedirect)
	}
}

// GetURLs - возвращает оригинальные ссылки по передаваемым сокращенным ссылкам.
func GetURLs(storage storage.Storage, baseURL string) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		// Аутентификация
		userUID := fmt.Sprintf("%s", req.Context().Value(UserKeyUID))

		var outputURLs []models.ResponseURL

		allURLs, err := storage.FindByUserUID(userUID)
		if err != nil {
			http.Error(res, "Error", http.StatusInternalServerError)
		}
		for _, url := range allURLs {
			outputURLs = append(outputURLs, models.ResponseURL{
				OriginalURL: url.OriginalURL, ShortURL: fmt.Sprintf("%s/%s", baseURL, url.ShortHash),
			})
		}

		res.Header().Set("Content-Type", "application/json")

		if len(outputURLs) == 0 {
			http.Error(res, "NoContent", http.StatusNoContent)
		} else {
			// Cериализуем ответ сервера
			enc := json.NewEncoder(res)
			if err := enc.Encode(outputURLs); err != nil {
				log.Printf("Error writing response: %s", err)
				return
			}
		}
	}
}

// DeleteURLs - обработчик удаления ссылок.
func DeleteURLs(mainStorage storage.Storage, inputCh chan []string) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		// Аутентификация
		userUID := fmt.Sprintf("%s", req.Context().Value(UserKeyUID))

		shortHashs, _ := io.ReadAll(req.Body)

		var shortsHashURL []string

		err := json.Unmarshal(shortHashs, &shortsHashURL)
		if err != nil {
			log.Printf("Error parsing JSON: %s", err)
			http.Error(res, "Bad JSON data", http.StatusBadRequest)
			return
		}

		// Удаление
		batches := chunkStrings(shortsHashURL, batchSize, userUID) // разбиваем на батчи массив коротких ссылок shortsHashURL - []string

		for _, batch := range batches {
			// Каждый батч удаляю в горутине
			inputCh <- batch
		}

		res.WriteHeader(http.StatusAccepted)
	}
}

func chunkStrings(arr []string, batchSize int, userUID string) [][]string {
	var batches [][]string

	// Проходим по массиву с шагом batchSize и добавляем подмассивы в batches
	for i := 0; i < len(arr); i += batchSize {
		end := i + batchSize

		// Убедимся, что индекс конца не выходит за границы массива
		if end > len(arr) {
			end = len(arr)
		}

		// Создаем новый батч с добавлением userUID
		batch := append([]string{userUID}, arr[i:end]...)

		// Добавляем подмассив в batches
		batches = append(batches, batch)
	}

	return batches
}