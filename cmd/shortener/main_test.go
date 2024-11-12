package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/PerfectStepCoder/shorturl/internal/handlers"
	"github.com/PerfectStepCoder/shorturl/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/PerfectStepCoder/shorturl/cmd/shortener/config"
	"github.com/go-resty/resty/v2"
)

const (
	testLengthShortURL = 10
	testBaseURL        = "http://localhost:8080"
)

func TestShorterURL(t *testing.T) {

	testCases := []struct {
		method       string
		body         string
		expectedCode int
		expectedBody string
	}{
		{method: http.MethodPost, body: "https://practicum.yandex.ru/", expectedCode: http.StatusCreated, expectedBody: "http://localhost:8080/42b3e75f92"},
		{method: http.MethodPost, body: "", expectedCode: http.StatusBadRequest, expectedBody: "URL not send\n"},
	}

	inMemoryStorage, _ := storage.NewStorageInMemory(testLengthShortURL)
	targetHandler := handlers.ShorterURL(inMemoryStorage, testBaseURL)

	srv := httptest.NewServer(targetHandler)
	defer srv.Close()

	for _, tc := range testCases {

		t.Run(tc.method, func(t *testing.T) {

			// Отправка HTTP-запроса
			req := resty.New().R()
			req.Method = tc.method
			req.SetBody(tc.body)
			req.URL = srv.URL
			resp, err := req.Send()

			if err != nil {
				t.Fatalf("Ошибка при отправке HTTP-запроса: %v", err)
			}

			// Проверка кода ответа
			if resp.StatusCode() != tc.expectedCode {
				t.Errorf("Код ответа не совпадает с ожидаемым. Ожидалось: %d, Получено: %d", tc.expectedCode, resp.StatusCode())
			}

			// Проверка тела ответа
			bodyString := string(resp.Body())
			if bodyString != tc.expectedBody {
				t.Errorf("Тело ответа не совпадает с ожидаемым. Ожидалось: %s, Получено: %s", tc.expectedBody, bodyString)
			}
		})
	}
}

func TestGetURLwithLoging(t *testing.T) {
	userUID := uuid.New().String()
	inMemoryStorage, _ := storage.NewStorageInMemory(testLengthShortURL)
	shortString, _ := inMemoryStorage.Save("https://yandex.ru/", userUID)
	assert.Equal(t, shortString, "77fca5950e")

	testCases := []struct {
		method       string
		body         string
		path         string
		expectedCode int
		expectedBody string
	}{
		{method: http.MethodGet, body: "", expectedCode: http.StatusOK, path: shortString, expectedBody: ""},
		{method: http.MethodGet, body: "", expectedCode: http.StatusNotFound, path: "NotExist", expectedBody: ""},
	}

	routes := chi.NewRouter()
	var logger, logFile = config.GetLogger()
	defer logFile.Close()
	routes.Get("/{id}", handlers.WithLogging(handlers.GetURL(inMemoryStorage), logger))
	srv := httptest.NewServer(routes)

	defer srv.Close()

	for _, tc := range testCases {

		t.Run(tc.method, func(t *testing.T) {

			// Отправка HTTP-запроса
			req := resty.New().R()
			req.Method = tc.method
			fmt.Print(srv.URL + "/" + tc.path + "\n") // вывод дополнительной информации
			resp, err := req.Get(srv.URL + "/" + tc.path)

			if err != nil {
				t.Fatalf("Ошибка при отправке HTTP-запроса: %v", err)
			}

			// Проверка кода ответа
			if resp.StatusCode() != tc.expectedCode {
				t.Errorf("Код ответа не совпадает с ожидаемым. Ожидалось: %d, Получено: %d", tc.expectedCode, resp.StatusCode())
			}
		})
	}
}

func TestObjectsURL(t *testing.T) {

	inMemoryStorage, _ := storage.NewStorageInMemory(testLengthShortURL)

	testCases := []struct {
		method       string
		body         string
		expectedCode int
		expectedBody string
	}{
		{method: http.MethodPost, body: "{\"url\":\"https://practicum.yandex.ru/\"}", expectedCode: http.StatusCreated, expectedBody: "{\"result\":\"http://localhost:8080/42b3e75f92\"}"},
		{method: http.MethodGet, body: "", expectedCode: http.StatusMethodNotAllowed, expectedBody: ""},
	}

	routes := chi.NewRouter()
	routes.Post("/api/shorten", handlers.ObjectShorterURL(inMemoryStorage, testBaseURL))
	srv := httptest.NewServer(routes)

	defer srv.Close()

	for _, tc := range testCases {

		t.Run(tc.method, func(t *testing.T) {

			// Отправка HTTP-запроса
			req := resty.New().R()
			req.Method = tc.method
			req.URL = srv.URL + "/api/shorten"

			if len(tc.body) > 0 {
				req.SetHeader("Content-Type", "application/json")
				req.SetBody(tc.body)
			}

			resp, err := req.Send()
			assert.NoError(t, err, "ошибка при отправке HTTP-запроса")

			// Проверка кода ответа
			if resp.StatusCode() != tc.expectedCode {
				t.Errorf("Код ответа не совпадает с ожидаемым. Ожидалось: %d, Получено: %d", tc.expectedCode, resp.StatusCode())
			}

			// Проверка содержимого тела запроса
			if tc.expectedBody != "" {
				assert.JSONEq(t, tc.expectedBody, string(resp.Body()), "содержимое ответа в формате JSON не совпадает")
			}
		})
	}
}

func TestGzipCompression(t *testing.T) {
	userUID := uuid.New().String()
	inMemoryStorage, _ := storage.NewStorageInMemory(testLengthShortURL)
	shortString, _ := inMemoryStorage.Save("https://practicum.yandex.ru/", userUID)
	assert.Equal(t, shortString, "42b3e75f92")

	testCases := []struct {
		method       string
		body         string
		contentType  string
		expectedCode int
		expectedBody string
		compress     bool
	}{
		{method: http.MethodPost, body: "{\"url\":\"https://yandex.ru/\"}", contentType: "application/json", compress: true, expectedCode: http.StatusCreated, expectedBody: "{\"result\":\"http://localhost:8080/77fca5950e\"}"},
		{method: http.MethodPost, body: "{\"url\":\"https://google.ru/\"}", contentType: "application/xml", compress: false, expectedCode: http.StatusCreated, expectedBody: "{\"result\":\"http://localhost:8080/41c9cc9cba\"}"},
	}

	routes := chi.NewRouter()
	routes.Post("/api/shorten", handlers.GzipCompress(handlers.ObjectShorterURL(inMemoryStorage, testBaseURL)))
	srv := httptest.NewServer(routes)

	defer srv.Close()

	for _, tc := range testCases {

		t.Run(tc.method, func(t *testing.T) {

			// Отправка HTTP-запроса
			req := resty.New().R()
			req.Method = tc.method
			req.URL = srv.URL + "/api/shorten"

			if len(tc.body) > 0 {
				req.SetHeader("Content-Type", tc.contentType)
				req.SetHeader("Accept-Encoding", "gzip")
				req.SetBody(tc.body)
			}

			resp, err := req.Send()
			assert.NoError(t, err, "ошибка при отправке HTTP-запроса")

			// Проверка кода ответа
			if resp.StatusCode() != tc.expectedCode {
				t.Errorf("Код ответа не совпадает с ожидаемым. Ожидалось: %d, Получено: %d", tc.expectedCode, resp.StatusCode())
			}

			// Проверка содержимого теля запроса
			if tc.expectedBody != "" {
				assert.JSONEq(t, tc.expectedBody, string(resp.Body()))
			}

			// Проверка сжатия данных
			if tc.compress {
				contentEncoding := resp.Header().Get("Content-Encoding")
				sendsGzip := strings.Contains(contentEncoding, "gzip")
				assert.True(t, sendsGzip)
			}
		})
	}
}

func TestAuthApiShorten(t *testing.T) {

	inMemoryStorage, _ := storage.NewStorageInMemory(testLengthShortURL)

	testCases := []struct {
		method       string
		path         string
		body         string
		contentType  string
		expectedCode int
		expectedBody string
		compress     bool
	}{
		{method: http.MethodPost, path: "/api/shorten", body: "{\"url\":\"https://yandex.ru/\"}", contentType: "application/json", compress: false, expectedCode: http.StatusCreated, expectedBody: "{\"result\":\"http://localhost:8080/77fca5950e\"}"},
		{method: http.MethodPost, path: "/api/shorten", body: "{\"url\":\"https://google.ru/\"}", contentType: "application/json", compress: false, expectedCode: http.StatusCreated, expectedBody: "{\"result\":\"http://localhost:8080/41c9cc9cba\"}"},
	}

	routes := chi.NewRouter()
	routes.Post("/api/shorten", handlers.ObjectShorterURL(inMemoryStorage, testBaseURL))
	routes.Get("/api/user/urls", handlers.CheckSignedCookie(handlers.Auth(handlers.GetURLs(inMemoryStorage, testBaseURL))))

	srv := httptest.NewServer(routes)
	defer srv.Close()

	for _, tc := range testCases {

		t.Run(tc.method, func(t *testing.T) {
			req := resty.New().R()
			req.Method = tc.method
			req.URL = srv.URL + tc.path

			// Установка куки
			resp, err := req.SetDebug(false).Send()
			assert.NoError(t, err, "ошибка при отправке HTTP-запроса")

			userUUID, found := findInCookie(resp)
			assert.True(t, found, fmt.Sprintf("Не найдена установленная кука: %s", resp.Cookies()))

			// Читаем записанные ссылки
			req = resty.New().R()
			req.URL = srv.URL + "/api/user/urls"
			req.Method = http.MethodGet
			cookie := &http.Cookie{
				Name:     "userUID",
				Value:    userUUID,
				Path:     "/",
				HttpOnly: true,  // Доступ только через HTTP
				Secure:   false, // Отправка только по HTTPS
			}
			req.SetCookie(cookie)
			_, err = req.SetDebug(false).Send()
			assert.NoError(t, err, "ошибка при отправке HTTP-запроса")
		})
	}
}

func TestBatchDelete(t *testing.T) {

	inMemoryStorage, _ := storage.NewStorageInMemory(testLengthShortURL)

	batch := "[{\"correlation_id\":\"8279bc80-2714-4767-8292-3e8328303e3f112\",\"original_url\":\"http://mail1.ru\"}, {\"correlation_id\":\"8d3f2ee8-af40-4c00-956b-da7415ba7e6e112\",\"original_url\":\"http://ya1.ru\"}]"

	resultBatch := "[{\"correlation_id\": \"8279bc80-2714-4767-8292-3e8328303e3f112\",\"short_url\": \"http://localhost:9999/8279bc80-2714-4767-8292-3e8328303e3f112\"},{\"correlation_id\": \"8d3f2ee8-af40-4c00-956b-da7415ba7e6e112\",\"short_url\": \"http://localhost:9999/8d3f2ee8-af40-4c00-956b-da7415ba7e6e112\"}]"

	testCases := []struct {
		method       string
		path         string
		body         string
		contentType  string
		expectedCode int
		expectedBody string
		compress     bool
	}{
		{method: http.MethodPost, path: "/api/shorten/batch", body: batch, contentType: "application/json", compress: false, expectedCode: http.StatusCreated, expectedBody: resultBatch},
	}

	inputCh := make(chan []string, 10000)

	routes := chi.NewRouter()
	routes.Post("/api/shorten/batch", handlers.ObjectsShorterURL(inMemoryStorage, testBaseURL))
	routes.Delete("/api/user/urls", handlers.Auth(handlers.DeleteURLs(mainStorage, inputCh)))
	srv := httptest.NewServer(routes)
	defer srv.Close()

	for _, tc := range testCases {

		t.Run(tc.method, func(t *testing.T) {
			req := resty.New().R()
			req.Method = tc.method
			req.URL = srv.URL + tc.path

			if len(tc.body) > 0 {
				req.SetHeader("Content-Type", tc.contentType)
				req.SetBody(tc.body)
			}

			_, err := req.Send()
			assert.NoError(t, err, "ошибка при отправке HTTP-запроса")

		})
	}

	// Удаление ссылок
	req := resty.New().R()
	req.Method = http.MethodDelete
	req.URL = srv.URL + "/api/user/urls"
	req.SetHeader("Content-Type", "application/json")
	req.SetBody("[\"8d3f2ee8-af40-4c00-956b-da7415ba7e6e112\", \"8279bc80-2714-4767-8292-3e8328303e3f112\"]")
	_, err := req.Send()
	assert.NoError(t, err, "ошибка при отправке HTTP-запроса")
	close(inputCh)
}

func TestPingDataBase(t *testing.T) {

	connectionStringDB := "http://localhost:5435/DB"
	routes := chi.NewRouter()
	routes.Get("/api/ping", handlers.PingDatabase(connectionStringDB))

	srv := httptest.NewServer(routes)
	defer srv.Close()

	req := resty.New().R()
	req.Method = http.MethodGet
	req.URL = srv.URL + "/api/ping"

	_, err := req.Send()
	assert.NoError(t, err, "ошибка при отправке HTTP-запроса")
}

func TestInitRoutes(t *testing.T) {

	var logger, logFile = config.GetLogger()
	defer logFile.Close()

	appSettings := config.ParseFlags()
	lengthInputCh := 1000
	inputCh := make(chan []string, lengthInputCh)
	mainStorage, _ = storage.NewStorageInMemory(lengthShortURL)
	routes := chi.NewRouter()

	err := initRoutes(routes, appSettings, logger, inputCh, mainStorage)
	assert.NoError(t, err)

}

func findInCookie(resp *resty.Response) (string, bool) { // userUUID, bool
	for _, cookie := range resp.Cookies() {
		fmt.Print(cookie.Name)
		if cookie.Name == "userUID" { // замените "userUID" на имя искомой cookie
			return cookie.Value, true
		}
	}
	return "", false
}
