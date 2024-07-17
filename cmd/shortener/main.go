package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PerfectStepCoder/shorturl/cmd/shortener/config"
	"github.com/PerfectStepCoder/shorturl/internal/handlers"
	"github.com/PerfectStepCoder/shorturl/internal/storage"
	"github.com/go-chi/chi/v5"
)

var inMemoryStorage *storage.StorageInMemory

const lengthShortURL = 10

func main() {

	appSettings := config.ParseFlags()

	inMemoryStorage = storage.NewStorage(lengthShortURL)
	routes := chi.NewRouter()

	routes.Post("/", handlers.ShorterURL(inMemoryStorage, appSettings.BaseURL))
	routes.Get("/{id}", handlers.GetURL(inMemoryStorage))

	fmt.Printf("Service is starting host: %s on port: %d\n", appSettings.ServiceNetAddress.Host,
		appSettings.ServiceNetAddress.Port)

	err := http.ListenAndServe(fmt.Sprintf(`%s:%d`, appSettings.ServiceNetAddress.Host,
		appSettings.ServiceNetAddress.Port), routes)

	if err != nil {
		log.Fatalf("error: %s", err)
	}
}
