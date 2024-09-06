// Package storage содержит функционал для персистентности данных
package storage

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// Storage - интерфейс для записи/чтения данных
type Storage interface {
	Save(value string, userUID string) (string, error)    // возвращает хеш ссылки
	Get(hashKey string) (string, bool)                    // возвращает origin ссылку или "" если не найдено
	Close()                                               // освобождение ресурсов
	FindByUserUID(userUID string) ([]ShortHashURL, error) // поиск сокращенных ссылок от пользователя
	IsDeleted(hashKey string) (bool, error)
	DeleteByUser(shortHashURL []string, userUID string) error
	//Cache(hashKey string)
}

type CorrelationURL struct {
	CorrelationID string
	OriginalURL   string
}

type ShortHashURL struct {
	ShortHash   string
	OriginalURL string
}

type CorrelationStorage interface {
	CorrelationSave(value string, correlationID string, userUID string) string           // возвращает хеш ссылки
	CorrelationGet(correlationID string) (string, bool)                                  // возвращает origin ссылку
	CorrelationsSave(correlationURLs []CorrelationURL, userUID string) ([]string, error) // возвращает срез хеш ссылок
}

// StorageFile - интерфейс для записи/чтения данных из файла
type StorageFile interface {
	LoadData(pathToFile string) int
	SaveData(pathToFile string) int
}

// Объединение интерфейсов
type PersistanceStorage interface {
	Storage
	StorageFile
	CorrelationStorage
}

func makeHash(value string, length int) string {
	output := ""
	hash := sha256.New()
	hash.Write([]byte(value))
	hashKey := hex.EncodeToString(hash.Sum(nil))
	output = hashKey[:length]
	return output
}

// TODO реализовать обертывание в эту ошибку все другие более "мелкие"
type StorageError struct {
	Err error
}

func NewStorageError(err error) error {
	return &StorageError{
		Err: err,
	}
}

func (se *StorageError) Error() string {
	return fmt.Sprintf("%v", se.Err)
}

type UniqURLError struct {
	ExistURL  string
	ShortHash string
}

func NewUniqURLError(existURL string, shortHash string) error {
	return &UniqURLError{existURL, shortHash}
}

func (e *UniqURLError) Error() string {
	return fmt.Sprintf("uniq error with %s", e.ExistURL)
}
