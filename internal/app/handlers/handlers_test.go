package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"io"
	"testing"
	"errors"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

type inputProvided struct {
	method   string
	path      string
	body     []byte
	store map[string]string
}

type outputDesired struct {
	code   int
	header map[string]string
	body   []byte
	store map[string]string
}

var tests = []struct {
	name string
	i    inputProvided
	o    outputDesired
}{
	{
		name: "Add new URL",
		i: inputProvided{
			method:   http.MethodPost,
			path:     "/",
			body:     []byte ("http://www.google.com"),
			store: map[string]string {},
		},
		o: outputDesired{
			code:   http.StatusCreated,
			header: nil,
			body:   nil,
			store: map[string]string { "": "http://www.google.com" },
		},
	},
	{
		name: "Try to get a non-existing URL #1",
		i: inputProvided{
			method:   http.MethodGet,
			path:     "/12345",
			body:     nil,
			store: map[string]string {},
		},
		o: outputDesired{
			code:   http.StatusBadRequest,
			header: nil,
			body:   nil,
			store: nil,
		},
	},
	{
		name: "Try to get a non-existing URL #2",
		i: inputProvided{
			method:   http.MethodGet,
			path:     "/12345",
			body:     nil,
			store: map[string]string { "54321": "http://www.google.com", },
		},
		o: outputDesired{
			code:   http.StatusBadRequest,
			header: nil,
			body:   nil,
			store: nil,
		},
	},
	{
		name: "Get an existing full URL #1",
		i: inputProvided{
			method:   http.MethodGet,
			path:     "/abc",
			body:     nil,
			store: map[string]string { "abc": "http://www.google.com", },
		},
		o: outputDesired{
			code:   http.StatusTemporaryRedirect,
			header: map[string]string { "Location": "http://www.google.com", },
			body:   nil,
			store: nil,
		},
	},
	{
		name: "Get an existing full URL #2",
		i: inputProvided{
			method:   http.MethodGet,
			path:     "/ABC",
			body:     nil,
			store: map[string]string { "abc": "http://www.google.com", "ABC": "http://www.yandex.ru", },
		},
		o: outputDesired{
			code:   http.StatusTemporaryRedirect,
			header: map[string]string { "Location": "http://www.yandex.ru", },
			body:   nil,
			store: nil,
		},
	},
}

/* This is a mock storage for test purposes using URLStorer interface. It implements
	the most simple approach with a memory-based map and a counter. The mock can also
	access the map directly without Add() / Get() for a faster test setup and checks.
*/
type urlStoreMock struct {
	i uint32
	s map[string]string
}
func (store *urlStoreMock) Add (url string) string {
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
func (store *urlStoreMock) Get (short string) (string, error) {
	url, ok := store.s[short]
	if !ok {
		return "", errors.New("URL does not exist in the store")
	}
	return url, nil
}

func TestSetRoute(t *testing.T) {
	for _, tt := range tests { 
		t.Run(tt.name, func(t *testing.T) {
			store := urlStoreMock { i: 0, s: tt.i.store, }
			router := chi.NewRouter()
			NewURLRouter("server:port", router, &store)
			server := httptest.NewServer(router)
			defer server.Close()
			
			request, err := http.NewRequest(tt.i.method, server.URL+tt.i.path, bytes.NewReader(tt.i.body))
			require.NoError(t, err)

			/* Provide CheckRedirect() to modify client's behavior for re-directs
			*/
			client := &http.Client {
				CheckRedirect: func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				},
			}

			response, err := client.Do(request)
			require.NoError(t, err)
			
			/* Check the response status code
			*/			
			if response.StatusCode != tt.o.code {
				t.Errorf("Expected status code %d, but got %d", tt.o.code, response.StatusCode)
			}

			/* Check the response header
			*/			
			for k, v := range tt.o.header {
				if r := response.Header.Get(k); r != v {
					t.Errorf("Expected header key \"%s\" = \"%s\", but key does not exist or = \"%s\"", k, v, r)
				}
			}

			/* Check the response body (if needed)
			*/
			if tt.o.body != nil {
				responseBody, err := io.ReadAll(response.Body)
				require.NoError(t, err)
				defer response.Body.Close()
				
				if !bytes.Equal(responseBody, tt.o.body) {
					t.Errorf("Expected body \"%s\", got \"%s\"", tt.o.body, responseBody)
				}
			}

			/* Check resulted URL store (if needed), but ignore map key values
			*/
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