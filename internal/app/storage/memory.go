package storage

import (
	"errors"
)

// Memory store impelementation. It uses a built-in map data type:
// - key is a short URL
// - value is the corresponding full URL

type URLStore struct {
	s map[string]string
}

func New() (*URLStore, error) {
	return &URLStore{
		s: map[string]string{},
	}, nil
}

func (store *URLStore) Add(url string) (string, error) {
	short := encode(url)
	if _,ok := store.s[short]; !ok {
		store.s[short] = url
	}
	return short, nil
}

func (store *URLStore) Get(short string) (string, error) {
	url, ok := store.s[short]
	if !ok {
		return "", errors.New(errNoURL)
	}
	return url, nil
}

func (store *URLStore) Close() error {
	return nil
}
