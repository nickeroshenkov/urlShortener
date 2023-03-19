package server

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/nickeroshenkov/urlShortener/internal/app/router"
	"github.com/nickeroshenkov/urlShortener/internal/app/storage"
)

type URLServer struct {
	http.Server
	Router *router.URLRouter
}

func New(cnf *Config) (*URLServer, error) {
	var srv URLServer

	sto, err := storage.New(*cnf.FileStoragePath)
	if err != nil {
		return nil, err
	}
	
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(router.DecompressRequest)
	r.Use(router.CompressResponse)

	srv.Router = router.New(*cnf.BaseURL, r, sto)

	srv.Addr = *cnf.ServerAddress
	srv.Handler = r

	return &srv, nil
}

func (srv *URLServer) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Server.Shutdown(ctx); err != nil {
		return err
	}
	if err := srv.Router.Close(); err != nil {
		return err
	}
	return nil
}
