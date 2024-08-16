package storage

import (
	"encoding/json"
	"os"
)

type ShortURL struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type Consumer struct {
	file    *os.File
	decoder *json.Decoder
}

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

func (c *Consumer) ReadShortURL() (*ShortURL, error) {
	shortURL := &ShortURL{}
	if err := c.decoder.Decode(shortURL); err != nil {
		return nil, err
	}
	return shortURL, nil
}

func (c *Consumer) Close() error {
	return c.file.Close()
}

type Producer struct {
	file    *os.File
	encoder *json.Encoder
}

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

func (p *Producer) WriteShortURL(shortURL *ShortURL) error {
	return p.encoder.Encode(&shortURL)
}

func (p *Producer) Close() error {
	return p.file.Close()
}
