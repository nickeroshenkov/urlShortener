package handlers

import (
	"fmt"
	"net/http"
	"io"

	"github.com/go-chi/chi/v5"

	"github.com/nickeroshenkov/urlShortener/internal/app/storage"
)

const (
	server = "localhost:8080"
)

func SetRoute(s storage.URLStorer, r chi.Router) {
	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		addURL(s, w, r)
	})
	r.Route("/{short}", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			getURL(s, w, r)
		})
	})
}

/* Both addURL() and getURL() take URLStorer as an argument. We could make them
	to be URLStorer methods instead, but it would make URLStorer scope less clear --
	"URLStorer" implies that the main scope is to store URLs, while HTTP handlers
	do things out of this scope.
*/

func addURL(s storage.URLStorer, w http.ResponseWriter, r *http.Request) {
	url, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	short := s.Add(string(url))
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, "http://" + server + "/" + short)
}

func getURL(s storage.URLStorer, w http.ResponseWriter, r *http.Request) {
	short := chi.URLParam(r, "short")
	if short == "" {
		http.Error(w, "Short URL identificator is missing", http.StatusBadRequest)
		return
	}
	url, err := s.Get(short)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}