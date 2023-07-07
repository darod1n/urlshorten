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
	SecretKey      string
}

func NewConfig() *Config {
	config := &Config{}
	config.Addr = os.Getenv("SERVER_ADDRESS")
	config.ServerHost = os.Getenv("BASE_URL")
	config.Path = os.Getenv("FILE_STORAGE_PATH")
	config.DataSourceName = os.Getenv("DATABASE_DSN")
	config.SecretKey = os.Getenv("SECRET_KEY")

	var flagAddr string
	var flagServerHost string
	var flagPath string
	var flagDataSourceName string
	var flagSecretKey string

	flag.StringVar(&flagAddr, "a", ":8080", "http-server address")
	flag.StringVar(&flagServerHost, "b", "http://localhost:8080", "base address result short url")
	flag.StringVar(&flagPath, "f", "", "File path")
	flag.StringVar(&flagDataSourceName, "d", "", "database data source name")
	flag.StringVar(&flagSecretKey, "s", "EnvSuperSecretKey", "secret key")
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

	if config.SecretKey == "" {
		config.SecretKey = flagSecretKey
	}

	return config
}
