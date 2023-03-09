package storage

import (
	"bufio"
	"io"
	"os"
	"strings"
)

// File store impelementation. It uses a text file with string pairs:
// - first string in a pair is a short URL
// - second string in a pair is the corresponding full URL
// - each string terminates with \n

type URLStoreFile struct {
	f  *os.File
	rw *bufio.ReadWriter
}

func NewURLStoreFile(filename string) (*URLStoreFile, error) {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0664)
	if err != nil {
		return nil, errOpen
	}
	rw := bufio.NewReadWriter(bufio.NewReader(f), bufio.NewWriter(f))
	return &URLStoreFile{
		f:  f,
		rw: rw,
	}, nil
}

func (store *URLStoreFile) Add(url string) (string, error) {
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
			return "", errRead
		}
		k = strings.TrimSpace(k)
		v = strings.TrimSpace(v)
		if v == url {
			return k, nil
		}
	}

	_, err1 := store.rw.WriteString(encode(url) + "\n")
	_, err2 := store.rw.WriteString(url + "\n")
	err3 := store.rw.Flush()
	if err1 != nil || err2 != nil || err3 != nil {
		return "", errWrite
	}

	return encode(url), nil
}

func (store *URLStoreFile) Get(short string) (string, error) {
	store.f.Seek(0, 0)
	store.rw.Reader.Reset(store.f)

	for {
		k, err1 := store.rw.ReadString('\n')
		v, err2 := store.rw.ReadString('\n')
		if len(k) == 0 && err1 == io.EOF {
			return "", errNoURL
		}
		if err1 != nil || err2 != nil {
			return "", errRead
		}
		k = strings.TrimSpace(k)
		v = strings.TrimSpace(v)
		if k == short {
			return v, nil
		}
	}
}

func (store *URLStoreFile) Close() error {
	err := store.f.Close()
	if err != nil {
		return errClose
	}
	return nil
}
