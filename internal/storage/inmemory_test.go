// Модуль содержит тесты интерфейса хранилища в памяти
package storage

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

const (
	testLengthShortURL = 10
)

// TestCreateURL - тестирование создание ссылки.
func TestCreateURL(t *testing.T) {

	inMemoryStorage, _ := NewStorageInMemory(testLengthShortURL)
	defer inMemoryStorage.Close()

	userUID := uuid.New().String()
	targetHash := "77fca5950e"
	shortString, _ := inMemoryStorage.Save("https://yandex.ru/", userUID)
	assert.Equal(t, shortString, targetHash)

	result, found := inMemoryStorage.Get(targetHash)
	assert.True(t, found)
	assert.Equal(t, "https://yandex.ru/", result)
}

// TestDeleteURL - тестирование удаление ссылки.
func TestDeleteURL(t *testing.T) {

	inMemoryStorage, _ := NewStorageInMemory(testLengthShortURL)
	defer inMemoryStorage.Close()

	userUID := uuid.New().String()
	shortString, _ := inMemoryStorage.Save("https://yandex.ru/", userUID)
	targetHash := "77fca5950e"
	assert.Equal(t, shortString, targetHash)

	inMemoryStorage.DeleteByUser([]string{targetHash}, userUID)
	_, found := inMemoryStorage.Get(targetHash)
	assert.False(t, found)

}

// TestFindURL - тестирование поиск ссылки.
func TestFindURL(t *testing.T) {

	inMemoryStorage, _ := NewStorageInMemory(testLengthShortURL)
	defer inMemoryStorage.Close()

	userUID := uuid.New().String()
	shortString, _ := inMemoryStorage.Save("https://yandex.ru/", userUID)
	assert.Equal(t, shortString, "77fca5950e")

	result, err := inMemoryStorage.FindByUserUID(userUID)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(result))
}

// TestCorrelationSaveGet - тестирование записи и чтения ссылок.
func TestCorrelationSaveGet(t *testing.T) {

	inMemoryStorage, _ := NewStorageInMemory(testLengthShortURL)
	defer inMemoryStorage.Close()

	userUID, correlationID := uuid.New().String(), uuid.New().String()

	correlationResult := inMemoryStorage.CorrelationSave("https://yandex.ru/", correlationID, userUID)
	assert.Equal(t, correlationID, correlationResult)

	result, found := inMemoryStorage.CorrelationGet(correlationID)
	assert.True(t, found)
	assert.Equal(t, fmt.Sprintf("%s|%s", "https://yandex.ru/", userUID), result)
}

// TestCorrelationsSaveGet - тесты записи и чтения ссылок массивами.
func TestCorrelationsSaveGet(t *testing.T) {

	inMemoryStorage, _ := NewStorageInMemory(testLengthShortURL)
	defer inMemoryStorage.Close()

	userUID := uuid.New().String()

	inputs := []CorrelationURL{
		{
			CorrelationID: uuid.New().String(),
			OriginalURL:   "https://yandex.ru/",
		},
		{
			CorrelationID: uuid.New().String(),
			OriginalURL:   "https://google.ru/",
		},
	}

	correlationResults, err := inMemoryStorage.CorrelationsSave(inputs, userUID)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(correlationResults))

}

// TestLoadSave - тесты записи и чтения хранилища в файле.
func TestLoadSave(t *testing.T) {

	inMemoryStorage, _ := NewStorageInMemory(testLengthShortURL)
	defer inMemoryStorage.Close()
	pathToFile := "noExist.db"

	countLoadRecords := inMemoryStorage.LoadData(pathToFile)
	assert.Equal(t, 0, countLoadRecords)

	countSaveRecords := inMemoryStorage.SaveData(pathToFile)
	assert.Equal(t, 0, countSaveRecords)
}
