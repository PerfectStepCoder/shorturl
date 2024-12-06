// Модуль config содержит настройки сервиса.
package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// NetAddress - хост на котором будет доступен сервис.
type NetAddress struct {
	Host string
	Port int
}

// String - функция печати хоста сервиса.
func (a NetAddress) String() string {
	return a.Host + ":" + strconv.Itoa(a.Port)
}

// Set - функция установки хоста из передоваемой строки
func (a *NetAddress) Set(s string) error {
	if a.Host == "" && a.Port == 0 {
		hp := strings.Split(s, ":")
		if len(hp) != 2 {
			return errors.New("need address in a form host:port")
		}
		port, err := strconv.Atoi(hp[1])
		if err != nil {
			return err
		}
		a.Host = hp[0]
		a.Port = port
	}
	return nil
}

// Settings - все настройки сервиса.
type Settings struct {
	ServiceNetAddress NetAddress
	BaseURL           string
	FileStoragePath   string
	DatabaseDSN       string
	ConfigNameFile    string
	SaveDBtoFile      bool
	AddProfileRoute   bool
	EnableTSL         bool
	TrustedSubnet     string
}

// Метод String для структуры Settings
func (s Settings) String() string {
	return fmt.Sprintf(
		"Settings:\n\tServiceNetAddress: %s\n\tBaseURL: %s\n\tFileStoragePath: %s\n\tDatabaseDSN: %s\n\tConfigNameFile: %s\n\tSaveDBtoFile: %v\n\tAddProfileRoute: %v\n\tEnableTSL: %v\n\tTrustedSubnet: %v",
		s.ServiceNetAddress, s.BaseURL, s.FileStoragePath, s.DatabaseDSN, s.ConfigNameFile, s.SaveDBtoFile, s.AddProfileRoute, s.EnableTSL, s.TrustedSubnet,
	)
}

// Config - структура для хранения данных из JSON
type ConfigJSON struct {
	ServerAddress   string `json:"server_address"`
	BaseURL         string `json:"base_url"`
	FileStoragePath string `json:"file_storage_path"`
	DatabaseDSN     string `json:"database_dsn"`
	EnableHTTPS     bool   `json:"enable_https"`
	TrustedSubnet   string `json:"trusted_subnet"`
}

// ParseConfig - функция для парсинга JSON-файла
func ParseJSONConfig(filename string) (*ConfigJSON, error) {
	// Открываем JSON-файл
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("не удалось открыть файл: %w", err)
	}
	defer file.Close()

	// Читаем содержимое файла
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать файл: %w", err)
	}

	// Парсим JSON в структуру Config
	var config ConfigJSON
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("ошибка при парсинге JSON: %w", err)
	}

	return &config, nil
}
