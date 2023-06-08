package config

import (
	"flag"
	"os"
)

type Config struct {
	Addr       string
	ServerHost string
	Path       string
}

func NewConfig() *Config {
	config := &Config{}
	config.Addr = os.Getenv("SERVER_ADDRESS")
	config.ServerHost = os.Getenv("BASE_URL")
	config.Path = os.Getenv("FILE_STORAGE_PATH")

	var flagAddr string
	var flagServerHost string
	var flagPath string

	flag.StringVar(&flagAddr, "a", ":8080", "http-server address")
	flag.StringVar(&flagServerHost, "b", "http://localhost:8080", "base address result short url")
	flag.StringVar(&flagPath, "f", "/tmp/short-url-db.json", "File path")
	flag.Parse()

	if config.Addr == "" {
		config.Addr = flagAddr
	}

	if config.ServerHost == "" {
		config.ServerHost = flagServerHost
	}

	if config.Path == "" {
		config.Path = flagPath
	}

	return config
}
