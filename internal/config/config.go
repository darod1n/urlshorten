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
	envAddr := os.Getenv("SERVER_ADDRESS")
	envServerHost := os.Getenv("BASE_URL")
	addr := flag.String("a", ":8080", "")
	serverHost := flag.String("b", "http://localhost:8080", "server adress")

	flag.Parse()

	if envAddr != "" {
		addr = &envAddr
	}

	if envServerHost != "" {
		serverHost = &envServerHost
	}

	return &Config{
		Addr:       *addr,
		ServerHost: *serverHost,
	}
}
