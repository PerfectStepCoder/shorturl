package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/PerfectStepCoder/shorturl/cmd/shortener/config"
	hdl "github.com/PerfectStepCoder/shorturl/internal/handlers"
	"github.com/PerfectStepCoder/shorturl/internal/storage"
	"github.com/go-chi/chi/v5"
)

var inMemoryStorage *storage.StorageInMemory

const lengthShortURL = 10

func main() {
	// Канал для получения сигналов
	sigs := make(chan os.Signal, 1)
	// Уведомлять о сигнале interrupt (Ctrl+C) и сигнале завершения
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	// Создаем канал для уведомления о завершении работы
	done := make(chan bool, 1)

	var logger, logFile = config.GetLogger()
	defer logFile.Close()

	appSettings := config.ParseFlags()
	fmt.Print("Settings:\n", appSettings, "\n")
	inMemoryStorage = storage.NewStorage(lengthShortURL)

	// Load
	loaded := inMemoryStorage.LoadData(appSettings.FileStoragePath)
	fmt.Printf("Loaded: %d recordes from file: %s\n", loaded, appSettings.FileStoragePath)

	routes := chi.NewRouter()
	routes.Post("/", hdl.GzipCompress(hdl.WithLogging(hdl.ShorterURL(inMemoryStorage, appSettings.BaseURL), logger)))
	routes.Get("/{id}", hdl.GzipCompress(hdl.WithLogging(hdl.GetURL(inMemoryStorage), logger)))
	routes.Post("/api/shorten", hdl.GzipCompress(hdl.ObjectShorterURL(inMemoryStorage, appSettings.BaseURL)))

	fmt.Printf("Service is starting host: %s on port: %d\n", appSettings.ServiceNetAddress.Host,
		appSettings.ServiceNetAddress.Port)

	go func() {
		err := http.ListenAndServe(fmt.Sprintf(`%s:%d`, appSettings.ServiceNetAddress.Host,
			appSettings.ServiceNetAddress.Port), routes)

		if err != nil {
			log.Fatalf("error: %s", err)
		}
	}()

	// Запуск горутины для обработки сигналов
	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Printf("Received signal: %s\n", sig)
		done <- true
	}()

	fmt.Println("Server is running...")
	fmt.Sprintf(`%s:%d`, appSettings.ServiceNetAddress.Host, appSettings.ServiceNetAddress.Port)

	// Ожидание сигнала завершения
	<-done
	fmt.Println("Shutting down server...")

	// Save
	saved := inMemoryStorage.SaveData(appSettings.FileStoragePath)
	fmt.Printf("Saved: %d recordes to file: %s\n", saved, appSettings.FileStoragePath)
}
