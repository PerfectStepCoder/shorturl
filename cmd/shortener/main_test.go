package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/go-chi/chi/v5"
	"github.com/PerfectStepCoder/shorturl/internal/handlers"
	"github.com/PerfectStepCoder/shorturl/internal/storage"
	"github.com/stretchr/testify/assert"

	"github.com/go-resty/resty/v2"
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
	targetHandler := handlers.ShorterURL(inMemoryStorage, "http://localhost:8080")
	
	srv := httptest.NewServer(targetHandler)
	defer srv.Close()
	
    for _, tc := range testCases {

        t.Run(tc.method, func(t *testing.T) {

			req := resty.New().R()
            req.Method = tc.method
			req.SetBody(tc.body)
            req.URL = srv.URL

            resp, err := req.Send()
            assert.NoError(t, err, "error making HTTP request")

            assert.Equal(t, tc.expectedCode, resp.StatusCode(), "Код ответа не совпадает с ожидаемым")

			assert.Equal(t, tc.expectedBody, string(resp.Body()), "Тело ответа не совпадает с ожидаемым")
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
        {method: http.MethodGet, body: "", expectedCode: http.StatusOK, path: shortString, expectedBody: ""},
		{method: http.MethodGet, body: "", expectedCode: http.StatusNotFound, path: "NotExist", expectedBody: ""},
    }

	routes := chi.NewRouter()
	routes.Get("/{id}", handlers.GetURL(inMemoryStorage))
	srv := httptest.NewServer(routes)

	defer srv.Close()

    for _, tc := range testCases {

        t.Run(tc.method, func(t *testing.T) {

			req := resty.New().R()
            req.Method = tc.method

			fmt.Print(srv.URL + "/" + tc.path + "\n")
			resp, err := req.Get(srv.URL + "/" + tc.path)
            assert.NoError(t, err, "error making HTTP request")

            assert.Equal(t, tc.expectedCode, resp.StatusCode(), "Код ответа не совпадает с ожидаемым")

        })
    }
}
