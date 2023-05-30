package storage

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
)

type Event struct {
	ID          int    `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type Events []Event

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

	if err := p.writer.WriteByte('\n'); err != nil {
		return err
	}

	return p.writer.Flush()
}

type Consumer struct {
	file    *os.File // файл для чтения
	scanner *bufio.Scanner
}

func NewConsumer(filename string) (*Consumer, error) {
	// открываем файл для чтения
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		file:    file,
		scanner: bufio.NewScanner(file),
	}, nil
}

func (c *Consumer) Close() error {
	return c.file.Close()
}

func (c *Consumer) ReadEvent() ([]Event, error) {
	var data []byte
	events := []Event{}
	event := Event{}

	for c.scanner.Scan() {
		data = c.scanner.Bytes()
		err := json.Unmarshal(data, &event)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil
}

func (c *Consumer) GetMap() map[string]string {
	em := make(map[string]string)
	events, err := c.ReadEvent()
	if err != nil {
		log.Println(err)
	}

	for _, event := range events {
		em[event.ShortURL] = event.OriginalURL
	}
	return em
}
