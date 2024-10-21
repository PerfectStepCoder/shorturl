package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	_ "net/http/pprof"
	"github.com/PerfectStepCoder/shorturl/cmd/shortener/config"
	hdl "github.com/PerfectStepCoder/shorturl/internal/handlers"
	"github.com/PerfectStepCoder/shorturl/internal/storage"
	"github.com/go-chi/chi/v5"
)

var mainStorage storage.PersistanceStorage

const lengthShortURL = 10

func main() {
	// Канал для получения сигналов
	sigs := make(chan os.Signal, 1)
	// Уведомлять о сигнале interrupt (Ctrl+C) и сигнале завершения
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	// Создаем канал для уведомления о завершении работы
	done := make(chan bool, 1)

	// Для обработки удаления urls
	inputCh := make(chan []string, 10000)  // TODO вынести 10000 в переменные окружения
		
	numWorkers := runtime.NumCPU()

	for i := 0; i < numWorkers; i++ {
		go func(inputCh chan []string) {
			for shortsHashURL := range inputCh {
				userUID := shortsHashURL[0]
				err := mainStorage.DeleteByUser(shortsHashURL[1:], userUID)
				if err != nil {
					log.Printf("Delete error: %s", err)
				}
			}
		}(inputCh)
	}

	var logger, logFile = config.GetLogger()
	defer logFile.Close()

	appSettings := config.ParseFlags()
	log.Print("Settings:\n", appSettings, "\n")
	log.Printf("Count core: %d\n", runtime.NumCPU())
	if appSettings.DatabaseDSN != "" {
		var err error
		mainStorage, err = storage.NewStorageInPostgres(appSettings.DatabaseDSN, lengthShortURL)
		if err != nil {
			log.Fatalf("Problem with database")
		}
	} else {
		mainStorage, _ = storage.NewStorageInMemory(lengthShortURL)
		// Load
		loaded := mainStorage.LoadData(appSettings.FileStoragePath)
		log.Printf("Loaded: %d recordes from file: %s\n", loaded, appSettings.FileStoragePath)
	}

	defer mainStorage.Close()

	routes := chi.NewRouter()

	// Middleware
	routes.Use(func(next http.Handler) http.Handler {
		return hdl.WithLogging(next.ServeHTTP, logger)
	})
	routes.Use(func(next http.Handler) http.Handler {
		return hdl.GzipCompress(next.ServeHTTP)
	})
	routes.Use(func(next http.Handler) http.Handler {
		return hdl.CheckSignedCookie(next.ServeHTTP)
	})
	
	// Регистрируем pprof маршруты
	routes.Mount("/debug/pprof/", http.DefaultServeMux)

	routes.Post("/", hdl.Auth(hdl.ShorterURL(mainStorage, appSettings.BaseURL)))
	routes.Get("/{id}", hdl.Auth(hdl.GetURL(mainStorage)))
	routes.Get("/api/user/urls", hdl.Auth(hdl.GetURLs(mainStorage, appSettings.BaseURL)))
	routes.Delete("/api/user/urls", hdl.Auth(hdl.DeleteURLs(mainStorage, inputCh)))
	routes.Post("/api/shorten", hdl.ObjectShorterURL(mainStorage, appSettings.BaseURL))
	routes.Post("/api/shorten/batch", hdl.ObjectsShorterURL(mainStorage, appSettings.BaseURL))
	routes.Get("/ping", hdl.PingDatabase(appSettings.DatabaseDSN))

	fmt.Printf("Service is starting host: %s on port: %d\n", appSettings.ServiceNetAddress.Host,
		appSettings.ServiceNetAddress.Port)

	go func() {
		err := http.ListenAndServe(fmt.Sprintf(`%s:%d`, appSettings.ServiceNetAddress.Host,
			appSettings.ServiceNetAddress.Port), routes)

		if err != nil {
			log.Printf("error: %s", err)
		}
	}()

	// Запуск горутины для обработки сигналов
	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Printf("Received signal: %s\n", sig)
		done <- true
	}()

	log.Printf("Server is running...")
	log.Printf(`%s:%d`, appSettings.ServiceNetAddress.Host, appSettings.ServiceNetAddress.Port)

	// Ожидание сигнала завершения
	<-done
	log.Println("Shutting down server...")
	close(inputCh)

	if appSettings.DatabaseDSN == "" {
		// Save
		saved := mainStorage.SaveData(appSettings.FileStoragePath)
		log.Printf("Saved: %d recordes to file: %s\n", saved, appSettings.FileStoragePath)
	}
}
