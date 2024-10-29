// Модуль для работы с флагами запуска сервиса.
package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Настройки по умолчанию.
const (
	baseURL         = "http://localhost:8080" // базовый адрес хоста сервиса
	fileStoragePath = "shorturls.data"        // путь к файлу хранилища ссылок
)

func splitHostPort(addr string) (string, int, error) {
	parts := strings.Split(addr, ":")
	if len(parts) != 2 {
		return "", 0, fmt.Errorf("invalid address format")
	}
	num, err := strconv.Atoi(parts[1])
	if err != nil {
		fmt.Println("can not str into number:", err)
		return parts[0], 0, nil
	}
	return parts[0], num, nil
}

// ParseFlags - функция для парсинга передаваемых флагов при старте сервиса.
func ParseFlags() Settings {
	appSettings := new(Settings)

	// Default
	appSettings.ServiceNetAddress.Host = ""
	appSettings.ServiceNetAddress.Port = 0

	flag.Var(&appSettings.ServiceNetAddress, "a", "Net address host:port")
	flag.StringVar(&appSettings.BaseURL, "b", baseURL, "Base url host:port")
	flag.StringVar(&appSettings.DatabaseDSN, "d", "", "DataBaseDSN connect to DB")
	flag.StringVar(&appSettings.FileStoragePath, "f", fileStoragePath, "Path to file of storage")
	flag.BoolVar(&appSettings.SaveDBtoFile, "s", false, "Save db to file")
	flag.BoolVar(&appSettings.AddProfileRoute, "p", false, "Add profiling route")
	flag.Parse()

	// Если есть переменные окружния они переписывают настройки
	if envBaseURL := os.Getenv("SHORTURL_BASE_URL"); envBaseURL != "" {
		appSettings.BaseURL = envBaseURL
	}
	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		appSettings.FileStoragePath = envFileStoragePath
	}
	if envDatabaseDSN := os.Getenv("SHORTURL_DATABASE_DSN"); envDatabaseDSN != "" {
		appSettings.DatabaseDSN = envDatabaseDSN
	}
	if envRunAddr := os.Getenv("SHORTURL_SERVER_ADDRESS"); envRunAddr != "" {
		host, port, err := splitHostPort(envRunAddr)
		if err == nil {
			appSettings.ServiceNetAddress.Host = host
			appSettings.ServiceNetAddress.Port = port
		}
	}

	if appSettings.ServiceNetAddress.Host == "" && appSettings.ServiceNetAddress.Port == 0 {
		appSettings.ServiceNetAddress.Host = "localhost"
		appSettings.ServiceNetAddress.Port = 8080
	}
	return *appSettings
}
