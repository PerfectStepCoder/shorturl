// Package storage содержит функционал для персистентности данных
package storage

import (
	"crypto/sha256"
	"encoding/hex"
)

// Storage - интерфейс для записи/чтения данных
type Storage interface {
	Save(value string) string
	Get(hashKey string) (string, bool)
	Close() // освобождение ресурсов
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
}

func make_hash(value string, length int) string {
	output := ""
	hash := sha256.New()
	hash.Write([]byte(value))
	hashKey := hex.EncodeToString(hash.Sum(nil))
	output = hashKey[:length]
	return output
}
