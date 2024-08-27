package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	// "net/http/cookiejar"
	// "net/url"
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
		}else {
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
		res.Header().Set("Location", originURL)
		res.WriteHeader(http.StatusTemporaryRedirect)
	}
}


func GetURLs(storage storage.Storage, baseURL string) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		// Вывод всех заголовков запроса
		for name, values := range req.Header {
			// Заголовок может иметь несколько значений, поэтому выводим их все
			for _, value := range values {
				log.Printf("%s: %s\n", name, value)
			}
		}

		// u, _ := url.Parse("https://httpbin.org")
		// cookies := req.Jar.Cookies(u)
		// fmt.Println("Cookies after first request:")
		// for _, cookie := range cookies {
		// 	fmt.Printf("%s = %s\n", cookie.Name, cookie.Value)
		// }

		// Аутентификация
		var userUID string
		cookies, err := req.Cookie("userUID")
		if err != nil {
			log.Print("No cookies")
			res.WriteHeader(http.StatusUnauthorized)
			return
		}else {
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
				log.Println("Error writing response:", err)
				return
			}
		}
	}
}
