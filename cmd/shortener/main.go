package main

import (
	"fmt"
	"net/http"

	"github.com/PerfectStepCoder/shorturl/internal/handlers"
	"github.com/PerfectStepCoder/shorturl/internal/storage"
	"github.com/go-chi/chi/v5"
)

var inMemoryStorage *storage.Storage

func main() {

	inMemoryStorage = storage.NewStorage(10);
	routes := chi.NewRouter()

	routes.Post("/", handlers.ShorterURL(inMemoryStorage))
    routes.Get("/{id}", handlers.GetURL(inMemoryStorage))

	PORT := 8080
	fmt.Printf("Service is starting ... on port: %d\n", PORT)
	err := http.ListenAndServe(fmt.Sprintf(`:%d`, PORT), routes)
	if err != nil {
		panic(err)
	}
}
