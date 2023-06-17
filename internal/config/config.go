package config

import (
	"flag"
	"os"
)

type Config struct {
	Addr           string
	ServerHost     string
	Path           string
	DataSourceName string
}

func NewConfig() *Config {
	config := &Config{}
	config.Addr = os.Getenv("SERVER_ADDRESS")
	config.ServerHost = os.Getenv("BASE_URL")
	config.Path = os.Getenv("FILE_STORAGE_PATH")
	config.DataSourceName = os.Getenv("DATABASE_DSN")

	var flagAddr string
	var flagServerHost string
	var flagPath string
	var flagDataSourceName string

	flag.StringVar(&flagAddr, "a", ":8080", "http-server address")
	flag.StringVar(&flagServerHost, "b", "http://localhost:8080", "base address result short url")
	flag.StringVar(&flagPath, "f", "", "File path")
	flag.StringVar(&flagDataSourceName, "d", "", "database data source name")
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

	if config.DataSourceName == "" {
		config.DataSourceName = flagDataSourceName
	}

	return config
}
