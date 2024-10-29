// Модуль содержит тесты бенчмарки
package alltests

import (
	"crypto/rand"
	"fmt"
	"github.com/google/uuid"
	"math/big"
	"testing"

	"github.com/PerfectStepCoder/shorturl/internal/storage"
	"github.com/stretchr/testify/assert"
)

// generateRandomLink генерирует случайную HTML-ссылку длиной lengthURL
func generateRandomLink(lengthURL int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	link := make([]byte, lengthURL)
	for i := range link {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			panic(err)
		}
		link[i] = letters[num.Int64()]
	}
	// Возвращаем HTML-ссылку
	return fmt.Sprintf(`https://example.com/%s`, string(link))
}

// generateUUID генерирует случайный UUID версии 4
func generateUUID() (string, error) {
	newUUID := uuid.New()
	return newUUID.String(), nil
}

// BenchmarkStorageInMemory - тестирует хранилище ссылок в памяти
func BenchmarkStorageInMemory(b *testing.B) {

	lengthShortURL := 20
	lengthURL := 10
	mainStorage, _ := storage.NewStorageInMemory(lengthShortURL)

	b.Run("storageInMemory", func(b *testing.B) {
		for i := 0; i < b.N; i++ {

			b.StopTimer() // останавливаем таймер
			rndUUID, _ := generateUUID()
			rndURL := generateRandomLink(lengthURL)
			b.StartTimer() // возобновляем таймер

			shortURL, err := mainStorage.Save(rndURL, rndUUID)
			assert.NoError(b, err)
			foundURL, found := mainStorage.Get(shortURL)
			assert.Equal(b, found, true)
			assert.Equal(b, foundURL, rndURL)
		}
	})
}
