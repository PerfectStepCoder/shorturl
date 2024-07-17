package storage

import (
	"crypto/sha256"
	"encoding/hex"
)

type StorageInMemory struct {
	data           map[string]string
	lengthShortURL int
}

func NewStorage(lengthShortURL int) *StorageInMemory {
	return &StorageInMemory{data: make(map[string]string), lengthShortURL: lengthShortURL}
}

func (s *StorageInMemory) Save(value string) string {
	hash := sha256.New()
	hash.Write([]byte(value))
	hashKey := hex.EncodeToString(hash.Sum(nil))
	hashKey = hashKey[:s.lengthShortURL]
	s.data[hashKey] = value
	return hashKey
}

func (s *StorageInMemory) Get(hashKey string) (string, bool) {
	value, exists := s.data[hashKey]
	return value, exists
}
