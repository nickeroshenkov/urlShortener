package storage

import ("errors")

type URLStorer interface {
	// Init ()
	Add (url string) int
	Get (id int) (string, error)
}

type URLStore struct {
	urls []string
}

// func (store *URLStore) Init () {
// 	store.urls = make([]string,0)
// }

func (store *URLStore) Add (url string) int {
	store.urls = append (store.urls, url)
	return len(store.urls)-1
}

func (store *URLStore) Get (id int) (string, error) {
	if id >= len(store.urls) {
		return "", errors.New ("URL does not exist in the store")
	}
	return store.urls[id], nil
}

