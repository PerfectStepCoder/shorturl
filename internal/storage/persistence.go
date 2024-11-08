// Модуль с реализацией сохранения данных в файл.
package storage

import (
	"encoding/json"
	"os"
)

// ShortURL - сохраняемая сущность в файл.
type ShortURL struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// Consumer - для работы с файлами.
type Consumer struct {
	file    *os.File
	decoder *json.Decoder
}

// NewConsumer - конструктор.
func NewConsumer(fileName string) (*Consumer, error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		file:    file,
		decoder: json.NewDecoder(file),
	}, nil
}

// ReadShortURL - чтение.
func (c *Consumer) ReadShortURL() (*ShortURL, error) {
	shortURL := &ShortURL{}
	if err := c.decoder.Decode(shortURL); err != nil {
		return nil, err
	}
	return shortURL, nil
}

// Close - освобождение ресурсов.
func (c *Consumer) Close() error {
	return c.file.Close()
}

// Producer - для работы с файлами.
type Producer struct {
	file    *os.File
	encoder *json.Encoder
}

// NewProducer - конструктор.
func NewProducer(fileName string) (*Producer, error) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &Producer{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

// WriteShortURL - запись
func (p *Producer) WriteShortURL(shortURL *ShortURL) error {
	return p.encoder.Encode(&shortURL)
}

// Close - освобождение ресурсов.
func (p *Producer) Close() error {
	return p.file.Close()
}
