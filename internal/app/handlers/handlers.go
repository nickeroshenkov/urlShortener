package handlers

import (
	"fmt"
	"github.com/nickeroshenkov/urlShortener/internal/app/storage"
	"io"
	"net/http"
	"strconv"
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

func Shortener(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		/* if len(r.URL.Query()) == 0 {
			fmt.Fprint(w, urlInputForm)
			return
		} */
		id_string := r.URL.Query().Get("id")
		if id_string == "" {
			http.Error(w, "Short URL identificator is missing", http.StatusBadRequest)
			return
		}
		id, err := strconv.ParseUint(id_string, 10, 0)
		if err != nil {
			http.Error(w, "Short URL identificator must be an unsigned integer", http.StatusBadRequest)
			return
		}
		if int(id) >= len(storage.UrlStore) {
			http.Error(w, "Short URL does not exist", http.StatusBadRequest)
			return
		}
		w.Header().Set("Location", storage.UrlStore[id])
		w.WriteHeader(http.StatusTemporaryRedirect)
	case http.MethodPost:
		// url := r.FormValue("url") // Form is used
		url, err := io.ReadAll(r.Body) // Form is not used yet, just read from the body
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Сheck if url is a URL indeed?
		storage.UrlStore = append(storage.UrlStore, string(url)) // Need to guard this with mutex?
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, "http://localhost:8080/?id=", len(storage.UrlStore)-1)
	default:
		http.Error(w, "Only GET or POST requests are allowed", http.StatusMethodNotAllowed)
	}
}