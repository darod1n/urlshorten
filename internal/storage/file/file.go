package file

import (
	"bufio"
	"encoding/json"
	"os"
)

type event struct {
	ID          int    `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type producer struct {
	file   *os.File // файл для записи
	writer *bufio.Writer
}

func newProducer(filename string) (*producer, error) {
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

func (p *producer) WriteEvent(event *event) error {
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
}

func newConsumer(filename string) (*consumer, error) {
	// открываем файл для чтения
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &consumer{
		file:    file,
		scanner: bufio.NewScanner(file),
	}, nil
}

func (c *consumer) Close() error {
	return c.file.Close()
}

func (c *consumer) ReadEvent() ([]event, error) {
	var data []byte
	events := []event{}
	event := event{}

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

func (c *consumer) GetMap() (map[string]string, error) {
	em := make(map[string]string)

	events, err := c.ReadEvent()
	if err != nil {
		return nil, err
	}

	for _, event := range events {
		em[event.ShortURL] = event.OriginalURL
	}
	return em, nil
}
