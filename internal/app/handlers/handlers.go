package handlers

import (
	"encoding/json"
	"fmt"
	"io"
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
	ur.chiRouter.Post("/", func(w http.ResponseWriter, r *http.Request) {
		ur.addURL(w, r)
	})
	ur.chiRouter.Post("/api/shorten", func(w http.ResponseWriter, r *http.Request) {
		ur.addURLAPI(w, r)
	})
	ur.chiRouter.Get("/{short}", func(w http.ResponseWriter, r *http.Request) {
		ur.getURL(w, r)
	})

	return &ur
}

func (ur *URLRouter) Close() error {
	return ur.urlStorer.Close()
}

func (ur URLRouter) addURL(w http.ResponseWriter, r *http.Request) {
	url, err1 := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err1 != nil {
		http.Error(w, err1.Error(), http.StatusInternalServerError)
		return
	}
	short, err2 := ur.urlStorer.Add(string(url))
	if err2 != nil {
		http.Error(w, err2.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, ur.baseURL+"/"+short)
}

func (ur URLRouter) addURLAPI(w http.ResponseWriter, r *http.Request) {
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
	short, err := ur.urlStorer.Add(string(request.URL))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Result = fmt.Sprintf("%s/%s", ur.baseURL, short)
	w.Header().Set(headerContentType, "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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
