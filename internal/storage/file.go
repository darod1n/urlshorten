package storage

import (
	"encoding/json"
	"os"
)

type Event struct {
	ID          int    `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type Producer struct {
	file *os.File // файл для записи
}

func NewProducer(filename string) (*Producer, error) {
	// открываем файл для записи в конец
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return &Producer{file: file}, nil
}

func (p *Producer) Close() error {
	// закрываем файл
	return p.file.Close()
}

func (p *Producer) WriteEvent(event *Event) error {
	data, err := json.Marshal(&event)
	if err != nil {
		return err
	}
	// добавляем перенос строки
	data = append(data, '\n')

	_, err = p.file.Write(data)
	return err
}

type Consumer struct {
	file *os.File // файл для чтения
}

func NewConsumer(filename string) (*Consumer, error) {
	// открываем файл для чтения
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &Consumer{file: file}, nil
}

func (c *Consumer) Close() error {
	// закрываем файл
	return c.file.Close()
}

func (c *Consumer) GetMap() map[string]string {
	return map[string]string{}
}
