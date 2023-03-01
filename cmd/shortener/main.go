package main

import (
	"flag"
	"log"
	"os"

	"github.com/nickeroshenkov/urlShortener/internal/app/server"
)

const (
	serverAddress = "localhost:8080"
	baseURL       = "http://localhost:8080/"
)

func main() {
	// Get config from the environment variables
	//
	a := os.Getenv("SERVER_ADDRESS")
	b := os.Getenv("BASE_URL")
	f := os.Getenv("FILE_STORAGE_PATH")

	// Get config from the flags
	//
	ap := flag.String("a", serverAddress, "specify server address in the form server:port")
	bp := flag.String("b", baseURL, "specify base URL in the form http://server:port/")
	fp := flag.String("f", "", "specify file storage path, empty one forces to use memory storage")
	flag.Parse()

	// Prioritize environment variables over flags
	//
	if a == "" {
		a = *ap
	}
	if b == "" {
		b = *bp
	}
	if f == "" {
		f = *fp
	}

	if err := server.Run(a, b, f); err != nil {
		log.Fatal(err)
	}
}
