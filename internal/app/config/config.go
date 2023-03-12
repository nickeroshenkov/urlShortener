package config

import (
	"os"
	"flag"
)

const (
	serverAddress = "localhost:8080"
	baseURL       = "http://localhost:8080"
)

type ServerConfig struct {
	ServerAddress   string
	BaseURL         string
	FileStoragePath string
}

// Load config from the environment variables, if values are non-empty
//
func (c *ServerConfig) LoadEVarsConditional() {
	a := os.Getenv("SERVER_ADDRESS")
	b := os.Getenv("BASE_URL")
	f := os.Getenv("FILE_STORAGE_PATH")
	if a != "" {
		c.ServerAddress = a
	}
	if b != "" {
		c.BaseURL = b
	}
	if f != "" {
		c.FileStoragePath = f
	}
}

// Load config from the flags, if values are non-empty
//
func (c *ServerConfig) LoadFlagsConditional() {
	ap := flag.String("a", serverAddress, "specify server address in the form server:port")
	bp := flag.String("b", baseURL, "specify base URL in the form http://server:port")
	fp := flag.String("f", "", "specify file storage path, empty one forces to use memory storage")
	flag.Parse()
	if *ap != "" {
		c.ServerAddress = *ap
	}
	if *bp != "" {
		c.BaseURL = *bp
	}
	if *fp != "" {
		c.FileStoragePath = *fp
	}
}

// Set the default config
//
func (c *ServerConfig) SetDefaults() {
	c.ServerAddress = serverAddress
	c.BaseURL = baseURL
	c.FileStoragePath = ""
}