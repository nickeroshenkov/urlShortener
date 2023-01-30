package main

import (
	"handlers"
	"net/http"
	"storage"
)

func main() {
	storage.UrlStore = make([]string, 0)
	http.HandleFunc("/", handlers.Shortener)
	http.ListenAndServe("localhost:8080", nil)
	// Consider to use log.Fatal(http.ListenAndServe("localhost:8080", nil)) instead
}
