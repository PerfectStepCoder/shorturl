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
	"github.com/stretchr/testify/assert"

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

func TestGetURL(t *testing.T) {

	inMemoryStorage, _ := storage.NewStorageInMemory(testLengthShortURL)
	shortString, _ := inMemoryStorage.Save("https://yandex.ru/")
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
	routes.Get("/{id}", handlers.GetURL(inMemoryStorage))
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

	inMemoryStorage, _ := storage.NewStorageInMemory(testLengthShortURL)
	shortString, _ := inMemoryStorage.Save("https://practicum.yandex.ru/")
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
