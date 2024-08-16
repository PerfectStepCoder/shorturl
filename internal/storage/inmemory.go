package storage

import (
	"io"
	"log"
	"sync"
)

type StorageInMemory struct {
	mu             sync.Mutex // синхронизация доступа к хранилищу
	data           map[string]string
	lengthShortURL int
}

func NewStorageInMemory(lengthShortURL int) (*StorageInMemory, error) {
	return &StorageInMemory{data: make(map[string]string), lengthShortURL: lengthShortURL}, nil
}

func (s *StorageInMemory) Save(value string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	hashKey := makeHash(value, s.lengthShortURL)
	// Проверка наличии ключа в map
	_, exists := s.data[hashKey]
	if exists {
		return hashKey, NewUniqURLError(value)
	} else {
		s.data[hashKey] = value
		return hashKey, nil
	}
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
		log.Print(err)
	}
	defer consumer.Close()
	for {
		shortURL, err := consumer.ReadShortURL()
		if err != nil {
			if err == io.EOF {
				break
			}
		}
		s.data[shortURL.ShortURL] = shortURL.OriginalURL
		count += 1
	}
	return count
}

func (s *StorageInMemory) SaveData(pathToFile string) int {
	producer, err := NewProducer(pathToFile)
	count := 0
	if err != nil {
		log.Print(err)
	}
	defer producer.Close()

	for shortURL, originURL := range s.data {
		newShortURL := ShortURL{
			UUID: shortURL, OriginalURL: originURL, ShortURL: shortURL,
		}
		if err := producer.WriteShortURL(&newShortURL); err != nil {
			log.Print(err)
		}
		count += 1
	}
	return count
}

func (s *StorageInMemory) CorrelationSave(value string, correlationID string) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[correlationID] = value
	return correlationID
}

func (s *StorageInMemory) CorrelationGet(correlationID string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	value, exists := s.data[correlationID]
	return value, exists
}

func (s *StorageInMemory) CorrelationsSave(correlationURLs []CorrelationURL) []string {

	var output []string

	for _, value := range correlationURLs {
		output = append(output, value.CorrelationID)
		s.CorrelationSave(value.OriginalURL, value.CorrelationID)
	}

	return output
}

func (s *StorageInMemory) Close() {
	s.data = nil
}
