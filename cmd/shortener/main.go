package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
)

var urlInputForm = `
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
`

var urlStore []string

func Shortener(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if len(r.URL.Query()) == 0 {
			fmt.Fprint(w, urlInputForm)
			return
		}
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "Short URL identificator is missing", http.StatusBadRequest)
			return
		}
		id64, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			http.Error(w, "Short URL identificator must be an unsigned integer", http.StatusBadRequest)
			return
		}
		if id64 >= uint64(len(urlStore)) {
			http.Error(w, "Short URL does not exist", http.StatusBadRequest)
			return
		}
		w.Header().Set("Location", urlStore[id64])
		w.WriteHeader(http.StatusTemporaryRedirect)
	case http.MethodPost:
		// url := r.FormValue("url")
		url, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Сheck if url is a URL indeed?
		urlStore = append(urlStore, string(url)) // Need to guard this with mutex?
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, "http://localhost:8080/?id=", len(urlStore)-1)
	default:
		http.Error(w, "Only GET or POST requests are allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	urlStore = make([]string, 0)
	http.HandleFunc("/", Shortener)
	http.ListenAndServe("localhost:8080", nil)
	// Consider to use log.Fatal(http.ListenAndServe("localhost:8080", nil)) instead
}
