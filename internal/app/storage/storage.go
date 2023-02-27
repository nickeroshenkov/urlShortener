package storage

import (
	"bufio"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"hash/fnv"
	"io"
	"log"
	"os"
	"strings"
)

type URLStorer interface {
	Add(url string) string
	Get(short string) (string, error)
}

// File store with 32-bit FNV-1a hashes and Base64 encoding (URL safe)

type URLStoreFile struct {
	f *os.File
}

func NewURLStoreFile(path string) *URLStoreFile {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Fatal("error opening the file store")
	}
	return &URLStoreFile{
		f: f,
	}
}

func (store *URLStoreFile) Close() {
	err := store.f.Close()
	if err != nil {
		log.Fatal("error closing the file store")
	}
}

func (store *URLStoreFile) Add(url string) string {
	store.f.Seek(0, 0)
	rw := bufio.NewReadWriter(bufio.NewReader(store.f), bufio.NewWriter(store.f))

	// Check if URL has already been stored
	//
	for {
		k, err1 := rw.ReadString('\n')
		v, err2 := rw.ReadString('\n')
		if len(k) == 0 && err1 == io.EOF {
			break
		}
		if err1 != nil || err2 != nil {
			log.Fatal("error reading the file store")
		}
		k = strings.TrimSpace(k)
		v = strings.TrimSpace(v)
		if v == url {
			return k
		}
	}

	h := make([]byte, 4)
	binary.LittleEndian.PutUint32(h, hash(url))
	short := base64.URLEncoding.EncodeToString(h)

	_, err1 := rw.WriteString(short + "\n")
	_, err2 := rw.WriteString(url + "\n")
	err3 := rw.Flush()
	if err1 != nil || err2 != nil || err3 != nil {
		log.Fatal("error writing the file store")
	}

	return short
}

func (store *URLStoreFile) Get(short string) (string, error) {
	store.f.Seek(0, 0)
	rw := bufio.NewReadWriter(bufio.NewReader(store.f), bufio.NewWriter(store.f))
	getError := errors.New("URL does not exist in the store")

	for {
		k, err1 := rw.ReadString('\n')
		v, err2 := rw.ReadString('\n')
		if len(k) == 0 && err1 == io.EOF {
			return "", getError
		}
		if err1 != nil || err2 != nil {
			log.Fatal("error reading the file store")
		}
		k = strings.TrimSpace(k)
		v = strings.TrimSpace(v)
		if k == short {
			return v, nil
		}
	}
}

// Memory store with 32-bit FNV-1a hashes and Base64 encoding (URL safe)

type URLStore struct {
	s map[string]string
}

func NewURLStore() *URLStore {
	return &URLStore{
		s: map[string]string{},
	}
}

func (store *URLStore) Add(url string) string {
	// Check if URL has already been stored
	//
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
