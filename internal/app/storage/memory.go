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

func NewURLStore() (*URLStore, error) {
	return &URLStore{
		s: map[string]string{},
	}, nil
}

func (store *URLStore) Add(url string) (string, error) {
	// Check if URL has already been stored
	//
	for k, v := range store.s {
		if v == url {
			return k, nil
		}
	}

	store.s[encode(url)] = url
	return encode(url), nil
}

func (store *URLStore) Get(short string) (string, error) {
	var getError = errors.New("URL does not exist in the store")
	url, ok := store.s[short]
	if !ok {
		return "", getError
	}
	return url, nil
}

func (store *URLStore) Close() error {
	return nil
}
