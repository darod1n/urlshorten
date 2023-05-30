package storage

import (
	"bufio"
	"encoding/json"
	"os"
)

type Event struct {
	ID          int    `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type Producer struct {
	file   *os.File // файл для записи
	writer *bufio.Writer
}

func NewProducer(filename string) (*Producer, error) {
	// открываем файл для записи в конец
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return &Producer{
		file:   file,
		writer: bufio.NewWriter(file),
	}, nil
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

	if _, err := p.writer.Write(data); err != nil {
		return err
	}

	// добавляем перенос строки

	if err := p.writer.WriteByte('\n'); err != nil {
		return err
	}

	return p.writer.Flush()
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
