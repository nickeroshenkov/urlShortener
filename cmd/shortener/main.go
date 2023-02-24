package main

import (
	"log"
	"os"

	"github.com/nickeroshenkov/urlShortener/internal/app/server"
)

const (
	serverAddress = "localhost:8080"
	baseURL       = "http://localhost:8080/"
)

func main() {
	s := os.Getenv("SERVER_ADDRESS")
	b := os.Getenv("BASE_URL")
	if s == "" {
		s = serverAddress
	}
	if b == "" {
		b = baseURL
	}
	if err := server.Run(s, b); err != nil {
		log.Fatal(err)
	}
}
