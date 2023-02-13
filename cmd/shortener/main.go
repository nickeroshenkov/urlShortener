package main

import (
	"log"

	"github.com/nickeroshenkov/urlShortener/internal/app/server"
)

func main() {
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}