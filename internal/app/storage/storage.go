package storage

import (
	"errors"
)

type URLStorer interface {
	Add(url string) uint32
	Get(id uint32) (string, error)
}

/* Remember to check for datarace once placed in memory
*/
type URLStore struct {
	i uint32
	s map[uint32]string
}

func New() *URLStore {
	return &URLStore {
		i: 0,
		s: map[uint32]string {},
	}
}

func (store *URLStore) Add(url string) uint32 {
	for k, v := range store.s {
		if v == url {
			return k
		}
	}
	store.i++
	store.s[store.i] = url
	return store.i
}

func (store *URLStore) Get(key uint32) (string, error) {
	var getError = errors.New("URL does not exist in the store")
	url, ok := store.s[key]
	if !ok {
		return "", getError
	}
	return url, nil
}