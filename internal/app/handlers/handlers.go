package handlers

import (
	"fmt"
	"net/http"
	"strconv"
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
	r.Route("/{urlID}", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			getURL(s, w, r)
		})
	})
}

/* Note both addURL() and getURL() take URLStorer as an argument. We could make them
	to be URLStorer methods instead, but it would make URLStorer scope less clear --
	"URLStorer" implies that the main scope is to store URLs, while HTTP handlers
	implementation is out of it.
*/

func addURL(s storage.URLStorer, w http.ResponseWriter, r *http.Request) {
	url, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	id := s.Add(string(url))
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, "http://" + server + "/" + strconv.FormatUint(uint64(id), 10))
}

func getURL(s storage.URLStorer, w http.ResponseWriter, r *http.Request) {
	idString := chi.URLParam(r, "urlID")
	if idString == "" {
		http.Error(w, "Short URL identificator is missing", http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseUint(idString, 10, 0)
	if err != nil {
		http.Error(w, "Short URL identificator must be an unsigned integer", http.StatusBadRequest)
		return
	}
	url, err := s.Get(uint32(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}