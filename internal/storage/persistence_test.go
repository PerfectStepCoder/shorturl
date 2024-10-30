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