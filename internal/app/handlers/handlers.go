package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/nickeroshenkov/urlShortener/internal/app/storage"
)

const (
	headerLocation    = "Location"
	headerContentType = "Content-Type"
)

type URLRouter struct {
	baseURL   string
	chiRouter chi.Router
	urlStorer storage.URLStorer
}

func NewURLRouter(s string, c chi.Router, u storage.URLStorer) *URLRouter {
	ur := URLRouter{
		baseURL:   s,
		chiRouter: c,
		urlStorer: u,
	}
	ur.chiRouter.Post("/api/shorten", func(w http.ResponseWriter, r *http.Request) {
		ur.addURL(w, r)
	})
	ur.chiRouter.Get("/{short}", func(w http.ResponseWriter, r *http.Request) {
		ur.getURL(w, r)
	})

	return &ur
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
	w.Header().Set(headerLocation, url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (ur URLRouter) addURL(w http.ResponseWriter, r *http.Request) {
	var request struct {
		URL string `json:"url"`
	}
	var response struct {
		Result string `json:"result"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	short := ur.urlStorer.Add(string(request.URL))
	response.Result = ur.baseURL + short
	w.Header().Set(headerContentType, "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
