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
	Close()
}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

// File store impelementation with 32-bit FNV-1a hashes and Base64 encoding (URL safe)
// It uses a text file with string pairs, each string terminates with \n
// - first string in a pair is a short URL
// - second string in a pair is the corresponding full URL

type URLStoreFile struct {
	f  *os.File
	rw *bufio.ReadWriter
}

func NewURLStoreFile(path string) *URLStoreFile {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Fatal("error opening the file store")
	}
	rw := bufio.NewReadWriter(bufio.NewReader(f), bufio.NewWriter(f))
	return &URLStoreFile{
		f:  f,
		rw: rw,
	}
}

func (store *URLStoreFile) Add(url string) string {
	store.f.Seek(0, 0)
	store.rw.Reader.Reset(store.f)
	store.rw.Writer.Reset(store.f)

	// Check if URL has already been stored
	//
	for {
		k, err1 := store.rw.ReadString('\n')
		v, err2 := store.rw.ReadString('\n')
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

	_, err1 := store.rw.WriteString(short + "\n")
	_, err2 := store.rw.WriteString(url + "\n")
	err3 := store.rw.Flush()
	if err1 != nil || err2 != nil || err3 != nil {
		log.Fatal("error writing the file store")
	}

	return short
}

func (store *URLStoreFile) Get(short string) (string, error) {
	store.f.Seek(0, 0)
	store.rw.Reader.Reset(store.f)
	getError := errors.New("URL does not exist in the store")

	for {
		k, err1 := store.rw.ReadString('\n')
		v, err2 := store.rw.ReadString('\n')
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

func (store *URLStoreFile) Close() {
	err := store.f.Close()
	if err != nil {
		log.Fatal("error closing the file store")
	}
}

// File store impelementation with 32-bit FNV-1a hashes and Base64 encoding (URL safe)
// It uses a built-in map data type:
// - key is a short URL
// - value is the corresponding full URL

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

func (store *URLStore) Close() {
}
