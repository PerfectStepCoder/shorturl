package storage

import (
	"crypto/sha256"
	"encoding/hex"
	"sync"
)

type StorageInMemory struct {
	mu sync.Mutex // синхронизация доступа к хранилищу
	data           map[string]string
	lengthShortURL int
}

func NewStorage(lengthShortURL int) *StorageInMemory {
	return &StorageInMemory{data: make(map[string]string), lengthShortURL: lengthShortURL}
}

func (s *StorageInMemory) Save(value string) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	hash := sha256.New()
	hash.Write([]byte(value))
	hashKey := hex.EncodeToString(hash.Sum(nil))
	hashKey = hashKey[:s.lengthShortURL]
	s.data[hashKey] = value
	return hashKey
}

func (s *StorageInMemory) Get(hashKey string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	value, exists := s.data[hashKey]
	return value, exists
}
