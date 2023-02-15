package storage

import (
	"errors"
	"encoding/binary"
	"encoding/base64"
	"hash/fnv"
)

type URLStorer interface {
	Add(url string) string
	Get(short string) (string, error)
}

/* Remember to check for datarace once placed in the memory
*/
type URLStore struct {
	s map[string]string
}

func NewURLStore() *URLStore {
	return &URLStore {
		s: map[string]string {},
	}
}

/* Current implementation of Add() and Get() use 32-bit FNV-1a hashes with Base64 encoding (URL safe)
*/

func (store *URLStore) Add(url string) string {
    /* Check if URL has already been stored
	*/
	for k, v := range store.s {
		if v == url {
			return k
		}
	}
	
	h := make([]byte, 4)
    binary.LittleEndian.PutUint32(h, hash(url))
	short := base64.URLEncoding.EncodeToString(h)
	
	store.s[short] = url
	return short
}

func (store *URLStore) Get(short string) (string, error) {
	var getError = errors.New("URL does not exist in the store")
	url, ok := store.s[short]
	if !ok {
		return "", getError
	}
	return url, nil
}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}
