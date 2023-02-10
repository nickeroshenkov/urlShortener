package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/nickeroshenkov/urlShortener/internal/app/handlers"
	"github.com/nickeroshenkov/urlShortener/internal/app/storage"
)

func Run() {
	var s storage.URLStore
	// s.Init() can be here

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	handlers.SetRoute(&s, r)

	http.ListenAndServe("localhost:8080", r)
	// Consider to use log.Fatal(http.ListenAndServe("localhost:8080", nil)) instead
}

