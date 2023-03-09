package handlers

import (
	"bytes"
	"compress/gzip"
	"fmt"
)

// gzipCompress() / gzipDecompress() are used by handlers tests only. They are the simple
// implementation of gzip coding from []byte to []byte from Yandex.Practicum platform
// "as is" with few adjustments (changing flate>gzip, using panic instead of errors).

// Compress сжимает слайс байт.
func gzipCompress(data []byte) []byte {
	var b bytes.Buffer
	// создаём переменную w — в неё будут записываться входящие данные,
	// которые будут сжиматься и сохраняться в bytes.Buffer
	w, err := gzip.NewWriterLevel(&b, gzip.BestSpeed)
	if err != nil {
		panic(fmt.Errorf("failed init compress writer: %v", err))
	}
	// запись данных
	_, err = w.Write(data)
	if err != nil {
		panic(fmt.Errorf("failed write data to compress temporary buffer: %v", err))
	}
	// обязательно нужно вызвать метод Close() — в противном случае часть данных
	// может не записаться в буфер b; если нужно выгрузить все упакованные данные
	// в какой-то момент сжатия, используйте метод Flush()
	err = w.Close()
	if err != nil {
		panic (fmt.Errorf("failed compress data: %v", err))
	}
	// переменная b содержит сжатые данные
	return b.Bytes()
}

// Decompress распаковывает слайс байт.
func gzipDecompress(data []byte) []byte {
	// переменная r будет читать входящие данные и распаковывать их
	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		panic(fmt.Errorf("failed init compress reader: %v", err))
	}
	defer r.Close()

	var b bytes.Buffer
	// в переменную b записываются распакованные данные
	_, err = b.ReadFrom(r)
	if err != nil {
		panic(fmt.Errorf("failed decompress data: %v", err))
	}

	return b.Bytes()
}
