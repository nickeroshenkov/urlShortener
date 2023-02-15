package main

import (
	"log"

	"github.com/nickeroshenkov/urlShortener/internal/app/server"
)

const (
	serverBaseURL = "localhost:8080"
)

func main() {
	if err := server.Run(serverBaseURL); err != nil {
		log.Fatal(err)
	}
}
