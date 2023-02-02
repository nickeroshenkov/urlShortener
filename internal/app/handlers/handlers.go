package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/nickeroshenkov/urlShortener/internal/app/storage"
)

/* var urlInputForm = `
<html>
    <head>
    <title></title>
    </head>
    <body>
        <form method="post">
            <label>Enter URL to shorten: </label><input type="text" name="url">
            <input type="submit" value="OK">
        </form>
    </body>
</html>
` */

func Shortener(s storage.URLStorer, w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		/* if len(r.URL.Query()) == 0 {
			fmt.Fprint(w, urlInputForm)
			return
		} */
		idString := r.URL.Query().Get("id")
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
		if err!=nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusTemporaryRedirect)
	case http.MethodPost:
		// url := r.FormValue("url") // Form is used
		url, err := io.ReadAll(r.Body) // Form is not used yet, just read from the body
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Сheck if url is a URL indeed?
		id := s.Add(string(url))
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, "http://localhost:8080/?id=", id)
	default:
		http.Error(w, "Only GET or POST requests are allowed", http.StatusMethodNotAllowed)
	}
}
