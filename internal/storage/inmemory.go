package storage

import (
	"crypto/sha256"
	"encoding/hex"
)

type Storage struct {
	data map[string]string
	maxLength int
}

func NewStorage(maxLength int) *Storage {
	return &Storage{data: make(map[string]string), maxLength: maxLength}
}

func (s *Storage) Save(value string) string {
	hash := sha256.New()
	hash.Write([]byte(value))
	hashKey := hex.EncodeToString(hash.Sum(nil))
	hashKey = hashKey[:s.maxLength]
	s.data[hashKey] = value
	return hashKey
}

func (s *Storage) Get(hashKey string) (string, bool) {
	value, exists := s.data[hashKey]
	return value, exists
}
