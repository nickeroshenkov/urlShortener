package handlers

import (
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/nickeroshenkov/urlShortener/internal/app/storage"
)

const (
	HeaderLocation = "Location"
)

type URLRouter struct {
	serverBaseURL string
	chiRouter     chi.Router
	urlStorer     storage.URLStorer
}

func NewURLRouter(s string, c chi.Router, u storage.URLStorer) *URLRouter {
	ur := URLRouter{
		serverBaseURL: s,
		chiRouter:     c,
		urlStorer:     u,
	}
	ur.chiRouter.Post("/", func(w http.ResponseWriter, r *http.Request) {
		ur.addURL(w, r)
	})
	ur.chiRouter.Route("/{short}", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			ur.getURL(w, r)
		})
	})
	return &ur
}

func (ur URLRouter) addURL(w http.ResponseWriter, r *http.Request) {
	url, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	short := ur.urlStorer.Add(string(url))
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, "http://"+ur.serverBaseURL+"/"+short)
}

func (ur URLRouter) getURL(w http.ResponseWriter, r *http.Request) {
	short := chi.URLParam(r, "short")
	if short == "" {
		http.Error(w, "Short URL identificator is missing", http.StatusBadRequest)
		return
	}
	url, err := ur.urlStorer.Get(short)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set(HeaderLocation, url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
