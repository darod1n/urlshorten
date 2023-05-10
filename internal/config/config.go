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

	var flagAddr string
	var flagServerHost string

	flag.StringVar(&flagAddr, "a", ":8080", "http-server address")
	flag.StringVar(&flagServerHost, "b", "http://localhost:8080", "base address result short url")
	flag.Parse()

	if config.Addr == "" {
		config.Addr = flagAddr
	}

	if config.ServerHost == "" {
		config.ServerHost = flagServerHost
	}

	return config
}
