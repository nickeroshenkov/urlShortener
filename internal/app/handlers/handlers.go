package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"io"

	"github.com/go-chi/chi/v5"

	"github.com/nickeroshenkov/urlShortener/internal/app/storage"
)

func SetRoute(s storage.URLStorer, r chi.Router) {
	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		addURL(s, w, r)
	})
	// r.Get("/", provideForm)
	r.Route("/{urlID}", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			getURL(s, w, r)
		})
		// r.Delete("/", ...)
	})
}

/* var newForm = `
<html>
    <head>
    <title></title>
    </head>
    <body>
        <form method="post">
            <label>Enter full URL to shorten, e.g. http://www.google.com : </label><input type="text" name="url">
            <input type="submit" value="OK">
        </form>
    </body>
</html>
`

func provideForm(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, newForm)
} */

func addURL(s storage.URLStorer, w http.ResponseWriter, r *http.Request) {
	url, err := io.ReadAll(r.Body) // Form is not used yet, just read from the body
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// url := r.FormValue("url") // Form is used

	// Сheck here if url is a URL indeed? + not ""

	id := s.Add(string(url))
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, "http://localhost:8080/", id)
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
	url, err := s.Get(int(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
