// Модуль содержит тесты интерфейса хранилища в памяти
package storage

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

const (
	testLengthShortURL = 10
)

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
