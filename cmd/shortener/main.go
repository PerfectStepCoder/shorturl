package main

import (
	"fmt"
	"net/http"

	"github.com/PerfectStepCoder/shorturl/internal/handlers"
	"github.com/PerfectStepCoder/shorturl/internal/storage"
)

var inMemoryStorage *storage.Storage

func main() {

	inMemoryStorage = storage.NewStorage(10);

	mux := http.NewServeMux()
	mux.HandleFunc("POST /", handlers.ShorterURL(inMemoryStorage))
	mux.HandleFunc("GET /{id}", handlers.GetURL(inMemoryStorage))

	PORT := 8080
	fmt.Printf("Service is starting ... on port: %d\n", PORT)
	err := http.ListenAndServe(fmt.Sprintf(`:%d`, PORT), mux)
	if err != nil {
		panic(err)
	}
}
