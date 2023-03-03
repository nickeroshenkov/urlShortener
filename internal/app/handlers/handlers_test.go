package handlers

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

type inputProvided struct {
	method   string
	path     string
	body     []byte
	compress bool              // true to allow server to compress the response
	store    map[string]string // nil is not allowed -- always initialize
}

type outputDesired struct {
	code   int
	header map[string]string
	body   []byte            // nil means do not test the response body
	store  map[string]string // nil means do not test the resulting store
}

var tests = []struct {
	name string
	i    inputProvided
	o    outputDesired
}{
	{
		name: "Add new URL via API #1",
		i: inputProvided{
			method:   http.MethodPost,
			path:     "/api/shorten",
			body:     []byte("{\"url\":\"http://www.google.com\"}"),
			compress: false,
			store:    map[string]string{},
		},
		o: outputDesired{
			code:   http.StatusOK,
			header: map[string]string{"Content-Type": "application/json"},
			// Assume the storage mock is added for the 1st time
			body:  []byte("{\"result\":\"http://server:port/1\"}\n"),
			store: map[string]string{"1": "http://www.google.com"},
		},
	},
	{
		name: "Add new URL via API #2",
		i: inputProvided{
			method:   http.MethodPost,
			path:     "/api/shorten",
			body:     []byte("{\"url\":\"http://www.yandex.ru\"}"),
			compress: false,
			store:    map[string]string{"100": "http://www.google.com"},
		},
		o: outputDesired{
			code:   http.StatusOK,
			header: map[string]string{"Content-Type": "application/json"},
			// Assume the storage mock is added for the 1st time
			body:  []byte("{\"result\":\"http://server:port/1\"}\n"),
			store: map[string]string{"100": "http://www.google.com", "1": "http://www.yandex.ru"},
		},
	},
	{
		name: "Add new URL via API #2 (with compression)",
		i: inputProvided{
			method:   http.MethodPost,
			path:     "/api/shorten",
			body:     []byte("{\"url\":\"http://www.yandex.ru\"}"),
			compress: true,
			store:    map[string]string{"100": "http://www.google.com"},
		},
		o: outputDesired{
			code:   http.StatusOK,
			header: map[string]string{"Content-Type": "application/json"},
			// Assume the storage mock is added for the 1st time
			body:  []byte("{\"result\":\"http://server:port/1\"}\n"),
			store: map[string]string{"100": "http://www.google.com", "1": "http://www.yandex.ru"},
		},
	},
	{
		name: "Add an already existing URL via API",
		i: inputProvided{
			method:   http.MethodPost,
			path:     "/api/shorten",
			body:     []byte("{\"url\":\"http://www.google.com\"}"),
			compress: false,
			store:    map[string]string{"100": "http://www.google.com", "101": "http://www.yandex.ru"},
		},
		o: outputDesired{
			code:   http.StatusOK,
			header: map[string]string{"Content-Type": "application/json"},
			body:   []byte("{\"result\":\"http://server:port/100\"}\n"),
			store:  map[string]string{"100": "http://www.google.com", "101": "http://www.yandex.ru"},
		},
	},
	{
		name: "Try to get a non-existing URL #1",
		i: inputProvided{
			method:   http.MethodGet,
			path:     "/12345",
			body:     nil,
			compress: false,
			store:    map[string]string{},
		},
		o: outputDesired{
			code:   http.StatusBadRequest,
			header: nil,
			body:   nil,
			store:  nil,
		},
	},
	{
		name: "Try to get a non-existing URL #2",
		i: inputProvided{
			method:   http.MethodGet,
			path:     "/12345",
			body:     nil,
			compress: false,
			store:    map[string]string{"54321": "http://www.google.com"},
		},
		o: outputDesired{
			code:   http.StatusBadRequest,
			header: nil,
			body:   nil,
			store:  nil,
		},
	},
	{
		name: "Get full URL #1",
		i: inputProvided{
			method:   http.MethodGet,
			path:     "/abc",
			body:     nil,
			compress: false,
			store:    map[string]string{"abc": "http://www.google.com"},
		},
		o: outputDesired{
			code:   http.StatusTemporaryRedirect,
			header: map[string]string{"Location": "http://www.google.com"},
			body:   nil,
			store:  nil,
		},
	},
	{
		name: "Get full URL #2",
		i: inputProvided{
			method:   http.MethodGet,
			path:     "/ABC",
			body:     nil,
			compress: false,
			store:    map[string]string{"abc": "http://www.google.com", "ABC": "http://www.yandex.ru"},
		},
		o: outputDesired{
			code:   http.StatusTemporaryRedirect,
			header: map[string]string{"Location": "http://www.yandex.ru"},
			body:   nil,
			store:  nil,
		},
	},
}

// This is a mock storage for test purposes using URLStorer interface. It implements
// the most simple approach with a memory-based map and a counter. The mock can also
// access the map directly without Add() / Get() for a faster test setup and checks.
type urlStoreMock struct {
	i uint32
	s map[string]string
}

func (store *urlStoreMock) Add(url string) string {
	for k, v := range store.s {
		if v == url {
			return k
		}
	}
	store.i++
	short := strconv.FormatUint(uint64(store.i), 10)
	store.s[short] = url
	return short
}
func (store *urlStoreMock) Get(short string) (string, error) {
	url, ok := store.s[short]
	if !ok {
		return "", errors.New("URL does not exist in the store")
	}
	return url, nil
}
func (store *urlStoreMock) Close() {
}

func TestSetRoute(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := urlStoreMock{i: 0, s: tt.i.store}
			router := chi.NewRouter()

			router.Use(DecompressRequest) // For gzip compression testing
			router.Use(CompressResponse)  // For gzip compression testing

			NewURLRouter("http://server:port/", router, &store)
			server := httptest.NewServer(router)
			defer server.Close()

			request, err := http.NewRequest(tt.i.method, server.URL+tt.i.path, bytes.NewReader(tt.i.body))
			require.NoError(t, err)

			// Prepare to test gzip compression
			//
			if tt.i.compress == true {
				request.Header.Set("Accept-Encoding", "gzip")
			} else {
				request.Header.Del("Accept-Encoding")
			}

			// Provide CheckRedirect() to modify client's behavior for re-directs
			//
			client := &http.Client{
				CheckRedirect: func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				},
			}

			response, err := client.Do(request)
			require.NoError(t, err)

			// Check the response status code
			//
			if response.StatusCode != tt.o.code {
				t.Errorf("Expected status code %d, but got %d", tt.o.code, response.StatusCode)
			}

			// Check the response header
			//
			for k, v := range tt.o.header {
				if r := response.Header.Get(k); r != v {
					t.Errorf("Expected header key \"%s\" = \"%s\", but key does not exist or = \"%s\"", k, v, r)
				}
			}

			// Check the response body (if needed)
			//
			if tt.o.body != nil {
				responseBody, err := io.ReadAll(response.Body)
				require.NoError(t, err)
				defer response.Body.Close()

				// Prepare to test gzip compression
				//
				if tt.i.compress == true {
					tt.o.body = gzipCompress(tt.o.body)
				}

				if !bytes.Equal(responseBody, tt.o.body) {
					t.Errorf("Expected body \"%s\", got \"%s\"", tt.o.body, responseBody)
				}
			}

			// Check the resulted store (if needed), but ignore map key values
			//
			if tt.o.store != nil {
				if len(tt.o.store) != len(store.s) {
					t.Errorf("Expected URL store size does not match the resulted one")
				}
			outer:
				for _, v := range tt.o.store {
					for _, w := range store.s {
						if v == w {
							continue outer
						}
					}
					t.Errorf("Expected URL store value %s does not exist", v)
				}
			}
		})
	}
}
