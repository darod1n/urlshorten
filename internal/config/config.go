package config

import (
	"flag"
	"os"
)

type Config struct {
	Addr       string
	ServerHost string
}

func NewConfig() *Config {
	config := &Config{}
	config.Addr = os.Getenv("SERVER_ADDRESS")
	config.ServerHost = os.Getenv("BASE_URL")

	if config.Addr == "" {
		flag.StringVar(&config.Addr, "a", ":8080", "http-server address")
	}

	if config.ServerHost == "" {
		flag.StringVar(&config.ServerHost, "b", "http://localhost:8080", "base address result short url")
	}
	flag.Parse()

	return config
}
