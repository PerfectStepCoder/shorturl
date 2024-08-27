package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"github.com/PerfectStepCoder/shorturl/internal/models"
	"github.com/PerfectStepCoder/shorturl/internal/storage"
	"github.com/go-chi/chi/v5"
)

func ShorterURL(mainStorage storage.Storage, baseURL string) http.HandlerFunc {
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

func GetURLs(storage storage.Storage, baseURL string) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		// Аутентификация
		var userUID string
		cookies, err := req.Cookie("userUID")
		if err != nil {
			log.Print("No cookies")
			encodedUserUID := req.Header.Get("Authorization")
			var validErr bool
			userUID, validErr = ValidateUserUID(encodedUserUID)
			if !validErr {
				res.WriteHeader(http.StatusUnauthorized)
				return
			}
		} else {
			userUID, _ = ValidateUserUID(cookies.Value) // обработка исключения не требуется
		}

		var outputURLs []models.ResponseURL

		allURLs, err := storage.FindByUserUID(userUID)
		if err == nil {
			for _, url := range allURLs {
				outputURLs = append(outputURLs, models.ResponseURL{
					OriginalURL: url.OriginalURL, ShortURL: fmt.Sprintf("%s/%s", baseURL, url.ShortHash),
				})
			}
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

func chunkStrings(arr []string, batchSize int) [][]string {
    var batches [][]string

    // Проходим по массиву с шагом batchSize и добавляем подмассивы в batches
    for i := 0; i < len(arr); i += batchSize {
        end := i + batchSize

        // Убедимся, что индекс конца не выходит за границы массива
        if end > len(arr) {
            end = len(arr)
        }

        // Добавляем подмассив в batches
        batches = append(batches, arr[i:end])
    }

    return batches
}

func DeleteURLs(mainStorage storage.Storage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		// Аутентификация
		var userUID string
		cookies, err := req.Cookie("userUID")
		if err != nil {
			log.Print("No cookies")
			encodedUserUID := req.Header.Get("Authorization")
			var validErr bool
			userUID, validErr = ValidateUserUID(encodedUserUID)
			if !validErr {
				//res.WriteHeader(http.StatusUnauthorized)
				res.WriteHeader(http.StatusAccepted)
				return
			}
		} else {
			userUID, _ = ValidateUserUID(cookies.Value) // обработка исключения не требуется
		}

		shortHashs, _ := io.ReadAll(req.Body)
		defer func() {
			if err := req.Body.Close(); err != nil {
				log.Printf("Could not close response body: %s", err)
			}
		}()

		var shortsHashURL []string
	
		err = json.Unmarshal(shortHashs, &shortsHashURL)
		if err != nil {
			log.Printf("Error parsing JSON: %s", err)
		}

		// Удаление
		batchSize := 50  // указываем размер батча
		batches := chunkStrings(shortsHashURL, batchSize)  // разбиваем на батчи массив коротких ссылок shortsHashURL - []string
		inputCh := make(chan []string, len(batches))
		
		var wg sync.WaitGroup
		numWorkers := 20
	
		for i := 0; i < numWorkers; i++ {
			wg.Add(1)
			go func(inputCh chan []string, wg *sync.WaitGroup) {
				defer wg.Done()
				for shortsHashURL := range inputCh {
					err = mainStorage.DeleteByUser(shortsHashURL, userUID)
					if err != nil {
						log.Printf("Delete error: %s", err)
					}
				}
			}(inputCh, &wg)
		}

		for i, batch := range batches {
			log.Printf("Batch %d: %v\n", i+1, batch)
			// Каждый батч удаляю в горутине
			inputCh <- batch
		}

		close(inputCh)

		// Ожидаем завершения всех воркеров
		//wg.Wait()

		res.WriteHeader(http.StatusAccepted)
	}
}
