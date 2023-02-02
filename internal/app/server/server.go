package server

import (
	"github.com/nickeroshenkov/urlShortener/internal/app/handlers"
	"github.com/nickeroshenkov/urlShortener/internal/app/storage"
	"net/http"
)

func Run() {
	var s storage.URLStore
	// s.Init()
	http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		handlers.Shortener(&s, w, r)
	})
	http.ListenAndServe("localhost:8080", nil)
	// Consider to use log.Fatal(http.ListenAndServe("localhost:8080", nil)) instead
}