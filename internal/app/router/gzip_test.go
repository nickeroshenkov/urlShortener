package router

import (
	"bytes"
	"compress/gzip"
	"errors"
)

// gzipCompress() / gzipDecompress() are used by handlers tests only. They are the simple
// implementation of gzip coding from []byte to []byte from Yandex.Practicum platform
// "as is" with few adjustments (changing flate>gzip, using panic instead of errors).

// Compress сжимает слайс байт.
func gzipCompress(data []byte) ([]byte, error) {
	var b bytes.Buffer
	// создаём переменную w — в неё будут записываться входящие данные,
	// которые будут сжиматься и сохраняться в bytes.Buffer
	w, err := gzip.NewWriterLevel(&b, gzip.BestSpeed)
	if err != nil {
		return nil, errors.New("failed init compress writer")
	}
	// запись данных
	_, err = w.Write(data)
	if err != nil {
		return nil, errors.New("failed write data to compress temporary buffer")
	}
	// обязательно нужно вызвать метод Close() — в противном случае часть данных
	// может не записаться в буфер b; если нужно выгрузить все упакованные данные
	// в какой-то момент сжатия, используйте метод Flush()
	err = w.Close()
	if err != nil {
		return nil, errors.New("failed compress data")
	}
	// переменная b содержит сжатые данные
	return b.Bytes(), nil
}

// Decompress распаковывает слайс байт.
func gzipDecompress(data []byte) ([]byte, error) {
	// переменная r будет читать входящие данные и распаковывать их
	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, errors.New("failed init compress reader")
	}
	defer r.Close()

	var b bytes.Buffer
	// в переменную b записываются распакованные данные
	_, err = b.ReadFrom(r)
	if err != nil {
		return nil, errors.New("failed decompress data")
	}

	return b.Bytes(), nil
}
