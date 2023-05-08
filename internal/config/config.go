package config

import (
	"flag"
)

type Config struct {
	Addr       string
	ServerHost string
}

func NewConfig() *Config {
	addr := flag.String("a", ":8080", "")
	serverHost := flag.String("b", "http://localhost:8080", "server adress")
	flag.Parse()
	return &Config{
		Addr:       *addr,
		ServerHost: *serverHost,
	}
}
