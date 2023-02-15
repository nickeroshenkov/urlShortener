package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/nickeroshenkov/urlShortener/internal/app/handlers"
	"github.com/nickeroshenkov/urlShortener/internal/app/storage"
)

func Run(serverBaseURL string) (err error) {
	s := storage.NewURLStore()
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	handlers.NewURLRouter(serverBaseURL, r, s)

	err = http.ListenAndServe(serverBaseURL, r)
		
	return
}