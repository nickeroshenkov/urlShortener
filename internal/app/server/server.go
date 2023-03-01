package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/nickeroshenkov/urlShortener/internal/app/handlers"
	"github.com/nickeroshenkov/urlShortener/internal/app/storage"
)

const (
	storeFilename = "store.txt"
)

func Run(serverAddress, baseURL, fileStoragePath string) (err error) {
	var s storage.URLStorer
	if fileStoragePath != "" {
		s = storage.NewURLStoreFile(fileStoragePath + "/" + storeFilename)
	} else {
		s = storage.NewURLStore()
	}
	defer s.Close()

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	handlers.NewURLRouter(baseURL, r, s)

	err = http.ListenAndServe(serverAddress, r)

	return
}
