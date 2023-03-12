package storage

import (
	"encoding/base64"
	"encoding/binary"
	"hash/fnv"
)

type URLStorer interface {
	Add(url string) (string, error)
	Get(short string) (string, error)
	Close() error
}

const (
	errNoURL = "URL does not exist in the store"
	errOpen  = "error opening the file store"
	errRead  = "error reading the file store"
	errWrite = "error writing the file store"
	errClose = "error closing the file store"
)

// All implementations use 32-bit FNV-1a hashes and Base64 encoding (URL safe)

func encode(url string) string {
	b := make([]byte, 4)
	h := fnv.New32a()
	h.Write([]byte(url))
	binary.LittleEndian.PutUint32(b, h.Sum32())
	return base64.URLEncoding.EncodeToString(b)
}
