package main

import (
	"fmt"
	"net/http/httptest"

	"github.com/PerfectStepCoder/shorturl/internal/handlers"
	"github.com/PerfectStepCoder/shorturl/internal/storage"
	"github.com/go-resty/resty/v2"
)

const (
	exampleLengthShortURL = 10
	exampleBaseURL        = "http://localhost:8080"
)

func ExampleTestShorterURL() {

	inMemoryStorage, _ := storage.NewStorageInMemory(exampleLengthShortURL)
	targetHandler := handlers.ShorterURL(inMemoryStorage, exampleBaseURL)

	srv := httptest.NewServer(targetHandler)
	defer srv.Close()

	// Отправка HTTP-запроса
	req := resty.New().R()
	req.Method = "POST"
	req.SetBody("https://practicum.yandex.ru/")
	req.URL = srv.URL
	resp, _ := req.Send()

	fmt.Print(resp)
}
