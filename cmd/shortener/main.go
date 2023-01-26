package main

import (
	"fmt"
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
		fmt.Fprint(w, urlStore[id64])
	case http.MethodPost:
		url := r.FormValue("url")
		// Then check if url is URL indeed
		urlStore = append(urlStore, url)
		// fmt.Fprint(w, url)
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
