// Модуль config содержит настройки сервиса.
package config

import (
	"runtime"
	"github.com/sirupsen/logrus"
	"github.com/PerfectStepCoder/shorturl/internal/storage"
)

func SetupServerSettings(log *logrus.Logger, lengthShortURL int) (Settings, storage.PersistanceStorage, error){

	appSettings := ParseFlags()

	// mainStorage - хранилище для записи и чтения обработанных ссылок.
	var mainStorage storage.PersistanceStorage

	log.Print("\n", appSettings, "\n")
	log.Printf("Count core: %d\n", runtime.NumCPU())
	if appSettings.DatabaseDSN != "" {
		var err error
		mainStorage, err = storage.NewStorageInPostgres(appSettings.DatabaseDSN, lengthShortURL)
		if err != nil {
			log.Fatalf("Problem with database")
			return appSettings, mainStorage, err
		}
	} else {
		mainStorage, _ = storage.NewStorageInMemory(lengthShortURL)
		// Load
		loaded := mainStorage.LoadData(appSettings.FileStoragePath)
		log.Printf("Loaded: %d recordes from file: %s\n", loaded, appSettings.FileStoragePath)
	}
	return appSettings, mainStorage, nil
}