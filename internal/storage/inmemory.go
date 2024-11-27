// Модуль содержит реализацию интерфейса хранилища в памяти
package storage

import (
	"fmt"
	"io"
	"log"
	"strings"
	"sync"
)

// StorageInMemory - хранилище в памяти ПК.
type StorageInMemory struct {
	mu             sync.Mutex // синхронизация доступа к хранилищу
	data           map[string]string
	lengthShortURL int
}

// NewStorageInMemory - конструктор.
func NewStorageInMemory(lengthShortURL int) (*StorageInMemory, error) {
	return &StorageInMemory{data: make(map[string]string), lengthShortURL: lengthShortURL}, nil
}

// Save - сохранение новой ссылки.
func (s *StorageInMemory) Save(value string, userUID string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	hashKey := makeHash(value, s.lengthShortURL)
	// Проверка наличии ключа в map
	_, exists := s.data[hashKey]
	if exists {
		return hashKey, NewUniqURLError(value, hashKey)
	} else {
		s.data[hashKey] = fmt.Sprintf("%s|%s", value, userUID) // hash -> originURL | userUUID
		return hashKey, nil
	}
}

// Get - чтение ссылки.
func (s *StorageInMemory) Get(hashKey string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	value, exists := s.data[hashKey]
	parts := strings.Split(value, "|")
	return parts[0], exists
}

// FindByUserUID - поиск ссылок по пользовательскому UID.
func (s *StorageInMemory) FindByUserUID(userUID string) ([]ShortHashURL, error) {
	var output []ShortHashURL

	for shortHash, originURLwithUserUID := range s.data {
		if strings.HasSuffix(originURLwithUserUID, userUID) {
			parts := strings.Split(originURLwithUserUID, "|")
			output = append(output, ShortHashURL{ShortHash: shortHash, OriginalURL: parts[0]})
		}
	}

	return output, nil
}

// LoadData загрузка данных из файла
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

// SaveData сохранение данных в файл
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

// CorrelationSave - сохранение данных (ссылка и идентификатор)
func (s *StorageInMemory) CorrelationSave(value string, correlationID string, userUID string) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[correlationID] = fmt.Sprintf("%s|%s", value, userUID)
	return correlationID
}

// CorrelationGet - чтение данных (ссылка и идентификатор)
func (s *StorageInMemory) CorrelationGet(correlationID string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	value, exists := s.data[correlationID]
	return value, exists
}

// CorrelationSave - сохранение данных (ссылок и идентификатор)
func (s *StorageInMemory) CorrelationsSave(correlationURLs []CorrelationURL, userUID string) ([]string, error) {

	var output []string

	for _, value := range correlationURLs {
		output = append(output, value.CorrelationID)
		s.CorrelationSave(value.OriginalURL, value.CorrelationID, userUID)
	}

	return output, nil
}

// IsDeleted - удалена ли ссылка.
func (s *StorageInMemory) IsDeleted(hashKey string) (bool, error) {
	// TODO нет реализации
	return false, nil
}

// DeleteByUser - удалить ссылку по пользовательскому UUID
func (s *StorageInMemory) DeleteByUser(shortHashURL []string, userUID string) error {

	for _, hash := range shortHashURL {
		if value, exists := s.data[hash]; exists {
			parts := strings.Split(value, "|")
			if len(parts) == 2 && parts[1] == userUID {
				delete(s.data, hash) // Удаляем ключ
			}
		}
	}
	// TODO продумать какое исключение передавать
	return nil
}

// Close - освобождение ресурсов
func (s *StorageInMemory) Close() {
	s.data = nil
}

func (s *StorageInMemory) CountURLs() (int, error) {
	output := 0
	output = len(s.data)
	return output, nil
}

func (s *StorageInMemory) CountUsers() (int, error) {

	output := 0

	set := make(map[string]struct{})

	// Функция для добавления ключа в множество
	add := func(key string) bool {
		if _, exists := set[key]; exists {
			// Ключ уже существует в множестве
			return false
		}
		// Добавляем ключ в множество
		set[key] = struct{}{}
		return true
	}

	for _, originURLwithUserUID := range s.data {
		parts := strings.Split(originURLwithUserUID, "|")
		if len(parts) == 2 && add(parts[1]) {
			output += 1
		}
	}

	return output, nil
}
