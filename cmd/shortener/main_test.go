package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/PerfectStepCoder/shorturl/internal/handlers"
	"github.com/PerfectStepCoder/shorturl/internal/storage"
	"github.com/stretchr/testify/assert"
)

func TestShorterURL(t *testing.T) {

    testCases := []struct {
        method       string
		body string
        expectedCode int
        expectedBody string
    }{
        {method: http.MethodPost, body: "https://practicum.yandex.ru/", expectedCode: http.StatusCreated, expectedBody: "http://localhost:8080/42b3e75f92"},
		{method: http.MethodPost, body: "", expectedCode: http.StatusBadRequest, expectedBody: "URL not send\n"},
    }

	inMemoryStorage := storage.NewStorage(10);
	targetHandler := handlers.ShorterURL(inMemoryStorage)

    for _, tc := range testCases {

        t.Run(tc.method, func(t *testing.T) {

            r := httptest.NewRequest(tc.method, "/", strings.NewReader(tc.body))
            w := httptest.NewRecorder()
 
			targetHandler(w, r)

            assert.Equal(t, tc.expectedCode, w.Code, "Код ответа не совпадает с ожидаемым")

			assert.Equal(t, tc.expectedBody, w.Body.String(), "Тело ответа не совпадает с ожидаемым")
        })
    }
}

func TestGetURL(t *testing.T) {

	inMemoryStorage := storage.NewStorage(10);
	shortString := inMemoryStorage.Save("https://practicum.yandex.ru/")
	assert.Equal(t, shortString, "42b3e75f92")

    testCases := []struct {
        method       string
		body string
		path string
        expectedCode int
        expectedBody string
    }{
        {method: http.MethodGet, body: "", expectedCode: http.StatusTemporaryRedirect, path: "/" + shortString, expectedBody: ""},
		{method: http.MethodGet, body: "", expectedCode: http.StatusNotFound, path: "/" + "NotExist", expectedBody: ""},
    }

	targetHandler := handlers.GetURL(inMemoryStorage)

    for _, tc := range testCases {

        t.Run(tc.method, func(t *testing.T) {

            r := httptest.NewRequest(tc.method, tc.path, strings.NewReader(tc.body))
			r.SetPathValue("id", tc.path[1:])

            w := httptest.NewRecorder()
 
			targetHandler(w, r)

            assert.Equal(t, tc.expectedCode, w.Code, "Код ответа не совпадает с ожидаемым")

        })
    }
}
