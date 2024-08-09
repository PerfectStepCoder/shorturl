package storage

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log"
	"sync"
)

type StorageInMemory struct {
	mu             sync.Mutex // синхронизация доступа к хранилищу
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

func (s *StorageInMemory) LoadData(pathToFile string) int {
	count := 0
	consumer, err := NewConsumer(pathToFile)
	if err != nil {
		log.Fatal(err)
	}
	defer consumer.Close()
	for {
		shortURL, err := consumer.ReadShortURL()
		if err != nil {
			if err == io.EOF {
				break
			}
		}
		s.data[shortURL.OriginalURL] = shortURL.ShortURL
		count += 1
	}
	return count
}

func (s *StorageInMemory) SaveData(pathToFile string) int {
	producer, err := NewProducer(pathToFile)
	count := 0
	if err != nil {
		log.Fatal(err)
	}
	defer producer.Close()

	for originURL, shortURL := range s.data {
		newShortURL := ShortURL{
			UUID: uint(count), OriginalURL: originURL, ShortURL: shortURL,
		}
		if err := producer.WriteShortURL(&newShortURL); err != nil {
			log.Fatal(err)
		}
		count += 1
	}
	return count
}
