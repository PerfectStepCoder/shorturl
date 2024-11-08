package storage

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewConsumer - тест для конструктора Consumer.
func TestNewConsumer(t *testing.T) {
	// Создание временного файла.
	file, err := os.CreateTemp("", "testfile.json")
	if err != nil {
		t.Fatalf("Не удалось создать временный файл: %v", err)
	}
	defer os.Remove(file.Name()) // Удаляем файл после теста.

	// Создание нового Consumer.
	consumer, err := NewConsumer(file.Name())
	assert.NoError(t, err)
	assert.NotNil(t, consumer)
	assert.Equal(t, file.Name(), consumer.file.Name())
}

// TestReadShortURL - тест для метода ReadShortURL.
func TestReadShortURL(t *testing.T) {
	// Создание временного файла.
	file, err := os.CreateTemp("", "testfile.json")
	if err != nil {
		t.Fatalf("Не удалось создать временный файл: %v", err)
	}
	defer os.Remove(file.Name()) // Удаляем файл после теста.

	// Записываем тестовые данные в файл.
	testData := ShortURL{UUID: "uuid", ShortURL: "fghfghfh", OriginalURL: "https://example.com"}
	encoder := json.NewEncoder(file)
	if err := encoder.Encode(testData); err != nil {
		t.Fatalf("Не удалось закодировать данные в JSON: %v", err)
	}

	// Закрываем файл, чтобы он был доступен для чтения.
	file.Close()

	// Создание нового Consumer.
	consumer, err := NewConsumer(file.Name())
	assert.NoError(t, err)

	// Чтение короткой ссылки.
	result, err := consumer.ReadShortURL()
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, testData.OriginalURL, result.OriginalURL)
}

// TestNewProducer - тест для конструктора Producer.
func TestNewProducer(t *testing.T) {
	// Создание временного файла.
	file, err := os.CreateTemp("", "testfile.json")
	if err != nil {
		t.Fatalf("Не удалось создать временный файл: %v", err)
	}
	defer os.Remove(file.Name()) // Удаляем файл после теста.

	// Создание нового Producer.
	producer, err := NewProducer(file.Name())
	assert.NoError(t, err)
	assert.NotNil(t, producer)
	assert.Equal(t, file.Name(), producer.file.Name())
}

// TestWriteShortURL - тест для метода WriteShortURL.
func TestWriteShortURL(t *testing.T) {
	// Создание временного файла.
	file, err := os.CreateTemp("", "testfile.json")
	if err != nil {
		t.Fatalf("Не удалось создать временный файл: %v", err)
	}
	defer os.Remove(file.Name()) // Удаляем файл после теста.

	// Создание нового Producer.
	producer, err := NewProducer(file.Name())
	assert.NoError(t, err)

	// Запись тестовых данных в файл.
	testData := ShortURL{
		UUID:        "123e4567-e89b-12d3-a456-426614174000",
		ShortURL:    "short.ly/abc123",
		OriginalURL: "https://example.com",
	}
	err = producer.WriteShortURL(&testData)
	assert.NoError(t, err)

	// Закрываем producer, чтобы записать данные в файл.
	err = producer.Close()
	assert.NoError(t, err)

	// Чтение записанных данных из файла для проверки.
	file.Seek(0, 0) // Сбросить указатель файла на начало.
	var readData ShortURL
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&readData)
	assert.NoError(t, err)

	// Проверка, что данные совпадают.
	assert.Equal(t, testData.UUID, readData.UUID)
	assert.Equal(t, testData.ShortURL, readData.ShortURL)
	assert.Equal(t, testData.OriginalURL, readData.OriginalURL)
}

// TestClose - тест для метода Close.
func TestClose(t *testing.T) {
	// Создание временного файла.
	file, err := os.CreateTemp("", "testfile.json")
	if err != nil {
		t.Fatalf("Не удалось создать временный файл: %v", err)
	}
	defer os.Remove(file.Name()) // Удаляем файл после теста.

	// Создание нового Producer.
	producer, err := NewProducer(file.Name())
	assert.NoError(t, err)

	// Закрытие producer.
	err = producer.Close()
	assert.NoError(t, err)

}
