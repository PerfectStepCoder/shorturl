package alltests

import (
	"crypto/rand"
	"fmt"
	"io"
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
	uuid := make([]byte, 16)
	_, err := io.ReadFull(rand.Reader, uuid)
	if err != nil {
		return "", err
	}

	// Устанавливаем 4-ю версию UUID (UUID v4)
	uuid[6] = (uuid[6] & 0x0f) | 0x40 // Устанавливаем 4-й бит (0100) для версии
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // Устанавливаем 2 старших бита для variant

	// Возвращаем строку в стандартном формате UUID: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}

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
