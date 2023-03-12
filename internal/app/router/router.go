package router

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
	router chi.Router
	storer storage.URLStorer
}

func New(s string, c chi.Router, u storage.URLStorer) *URLRouter {
	rou := URLRouter{
		baseURL:   s,
		router: c,
		storer: u,
	}
	rou.router.Post("/", func(w http.ResponseWriter, r *http.Request) {
		rou.addURL(w, r)
	})
	rou.router.Post("/api/shorten", func(w http.ResponseWriter, r *http.Request) {
		rou.addURLAPI(w, r)
	})
	rou.router.Get("/{short}", func(w http.ResponseWriter, r *http.Request) {
		rou.getURL(w, r)
	})

	return &rou
}

func (rou *URLRouter) Close() error {
	return rou.storer.Close()
}

func (rou URLRouter) addURL(w http.ResponseWriter, r *http.Request) {
	url, err1 := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err1 != nil {
		http.Error(w, err1.Error(), http.StatusInternalServerError)
		return
	}
	short, err2 := rou.storer.Add(string(url))
	if err2 != nil {
		http.Error(w, err2.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, rou.baseURL+"/"+short)
}

func (rou URLRouter) addURLAPI(w http.ResponseWriter, r *http.Request) {
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
	short, err := rou.storer.Add(string(request.URL))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Result = fmt.Sprintf("%s/%s", rou.baseURL, short)
	w.Header().Set(headerContentType, "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (rou URLRouter) getURL(w http.ResponseWriter, r *http.Request) {
	short := chi.URLParam(r, "short")
	if short == "" {
		http.Error(w, "Short URL identificator is missing", http.StatusBadRequest)
		return
	}
	url, err := rou.storer.Get(short)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set(headerLocation, url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
