package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/nickeroshenkov/urlShortener/internal/app/config"	
	"github.com/nickeroshenkov/urlShortener/internal/app/handlers"
	"github.com/nickeroshenkov/urlShortener/internal/app/storage"
)

func Run(c *config.Config) error {
	var s storage.URLStorer
	var err error
	if c.FileStoragePath != "" {
		s, err = storage.NewURLStoreFile(c.FileStoragePath)
	} else {
		s, err = storage.NewURLStore()
	}
	if err != nil {
		return err
	}
	defer s.Close()

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(handlers.DecompressRequest)
	r.Use(handlers.CompressResponse)

	handlers.NewURLRouter(c.BaseURL, r, s)

	return http.ListenAndServe(c.ServerAddress, r)
}
