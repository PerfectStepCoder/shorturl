// Package storage содержит функционал для персистентности данных
package storage

// Storage - интерфейс для записи/чтения данных
type Storage interface {
	Save(value string) string
	Get(hashKey string) (string, bool)
}
