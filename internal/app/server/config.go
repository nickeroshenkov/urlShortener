package server

import (
	"os"
	"flag"
)

const (
	defaultServerAddress = "localhost:8080"
	defaultBaseURL       = "http://localhost:8080"
)

type Config struct {
	ServerAddress   *string
	BaseURL         *string
	FileStoragePath *string
}

func NewConfig() *Config {
	var c Config

	c.ServerAddress = flag.String("a", defaultServerAddress, "specify server address in the form server:port")
	c.BaseURL = flag.String("b", defaultBaseURL, "specify base URL in the form http://server:port")
	c.FileStoragePath = flag.String("f", "", "specify file storage path, empty one forces to use memory storage")

	return &c
}

func (c *Config) Parse () {
	flag.Parse()

	a := os.Getenv("SERVER_ADDRESS")
	b := os.Getenv("BASE_URL")
	f := os.Getenv("FILE_STORAGE_PATH")
	if a != "" {
		c.ServerAddress = &a
	}
	if b != "" {
		c.BaseURL = &b
	}
	if f != "" {
		c.FileStoragePath = &f
	}
}
