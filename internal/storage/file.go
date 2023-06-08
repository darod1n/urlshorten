package storage

import (
	"bufio"
	"encoding/json"
	"os"

	"go.uber.org/zap"
)

type Event struct {
	ID          int    `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type Events []Event

type producer struct {
	file   *os.File // файл для записи
	writer *bufio.Writer
}

func NewProducer(filename string) (*producer, error) {
	// открываем файл для записи в конец
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return &producer{
		file:   file,
		writer: bufio.NewWriter(file),
	}, nil
}

func (p *producer) Close() error {
	// закрываем файл
	return p.file.Close()
}

func (p *producer) WriteEvent(event *Event) error {
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

type consumer struct {
	file    *os.File // файл для чтения
	scanner *bufio.Scanner
	l       *zap.SugaredLogger
}

func NewConsumer(filename string, l *zap.SugaredLogger) (*consumer, error) {
	// открываем файл для чтения
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &consumer{
		file:    file,
		scanner: bufio.NewScanner(file),
		l:       l,
	}, nil
}

func (c *consumer) Close() error {
	return c.file.Close()
}

func (c *consumer) ReadEvent() ([]Event, error) {
	var data []byte
	events := []Event{}
	event := Event{}

	for c.scanner.Scan() {
		data = c.scanner.Bytes()
		err := json.Unmarshal(data, &event)
		if err != nil {
			c.l.Errorf("failed to unmarshal: %v", err)
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil
}

func (c *consumer) GetMap() map[string]string {
	em := make(map[string]string)
	events, err := c.ReadEvent()
	if err != nil {
		c.l.Errorf("failed to read event: %v", err)
	}

	for _, event := range events {
		em[event.ShortURL] = event.OriginalURL
	}
	return em
}
