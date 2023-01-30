package server

import (
	"github.com/nickeroshenkov/urlShortener/internal/app/handlers"
	"net/http"
)

func Run() {
	http.HandleFunc("/", handlers.Shortener)
	http.ListenAndServe("localhost:8080", nil)
	// Consider to use log.Fatal(http.ListenAndServe("localhost:8080", nil)) instead
}
