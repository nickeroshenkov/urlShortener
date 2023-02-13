package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"io"
	"testing"
	"errors"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type inputProvided struct {
	method   string
	path      string
	body     []byte
	URLStore map[uint32]string
}

type outputDesired struct {
	code   int
	header map[string]string
	body   []byte
	URLStore map[uint32]string
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
			URLStore: map[uint32]string {},
		},
		o: outputDesired{
			code:   http.StatusCreated,
			header: nil,
			body:   nil,
			URLStore: map[uint32]string { 1: "http://www.google.com", },
		},
	},
	{
		name: "Try to get with a wrong value",
		i: inputProvided{
			method:   http.MethodGet,
			path:     "/a",
			body:     nil,
			URLStore: map[uint32]string {},
		},
		o: outputDesired{
			code:   http.StatusBadRequest,
			header: nil,
			body:   nil,
			URLStore: map[uint32]string {},
		},
	},
	{
		name: "Try to get a non-existing URL #1",
		i: inputProvided{
			method:   http.MethodGet,
			path:     "/1",
			body:     nil,
			URLStore: map[uint32]string {},
		},
		o: outputDesired{
			code:   http.StatusBadRequest,
			header: nil,
			body:   nil,
			URLStore: map[uint32]string {},
		},
	},
	{
		name: "Try to get a non-existing URL #2",
		i: inputProvided{
			method:   http.MethodGet,
			path:     "/2",
			body:     nil,
			URLStore: map[uint32]string { 1: "http://www.google.com", },
		},
		o: outputDesired{
			code:   http.StatusBadRequest,
			header: nil,
			body:   nil,
			URLStore: map[uint32]string { 1: "http://www.google.com", },
		},
	},
	{
		name: "Get an existing full URL #1",
		i: inputProvided{
			method:   http.MethodGet,
			path:     "/1",
			body:     nil,
			URLStore: map[uint32]string { 1: "http://www.google.com", },
		},
		o: outputDesired{
			code:   http.StatusTemporaryRedirect,
			header: map[string]string { "Location": "http://www.google.com", },
			body:   nil,
			URLStore: map[uint32]string { 1: "http://www.google.com", },
		},
	},
	{
		name: "Get an existing full URL #2",
		i: inputProvided{
			method:   http.MethodGet,
			path:     "/2",
			body:     nil,
			URLStore: map[uint32]string { 1: "http://www.google.com", 2: "http://www.yandex.ru", },
		},
		o: outputDesired{
			code:   http.StatusTemporaryRedirect,
			header: map[string]string { "Location": "http://www.yandex.ru", },
			body:   nil,
			URLStore: map[uint32]string { 1: "http://www.google.com", 2: "http://www.yandex.ru", },
		},
	},
}

/* This is a mock using storage.URLStorer interface. It is implemented using a memory-based map.
	At the time being, storage.URLStore uses the same approach, but this can change in the future.
	Also, the mock allows a direct access to the map without Add()/Get() for the purpose of easier
	tests setup.
*/
type urlStoreMock struct {
	i uint32
	s map[uint32]string
}
func (store *urlStoreMock) Add (url string) uint32 {
	for k, v := range store.s {
		if v == url {
			return k
		}
	}
	store.i++
	store.s[store.i] = url
	return store.i
}
func (store *urlStoreMock) Get (key uint32) (string, error) {
	url, ok := store.s[key]
	if !ok {
		return "", errors.New("URL does not exist in the store")
	}
	return url, nil
}

func TestSetRoute(t *testing.T) {
	for _, tt := range tests { 
		t.Run(tt.name, func(t *testing.T) {
			store := urlStoreMock { i: 0, s: tt.i.URLStore, }
			
			router := chi.NewRouter()
			SetRoute (&store, router)
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
			
			/* Response status code check
			*/			
			if response.StatusCode != tt.o.code {
				t.Errorf("Expected status code %d, but got %d", tt.o.code, response.StatusCode)
			}

			/* Response header check
			*/			
			for k, v := range tt.o.header {
				if r := response.Header.Get(k); r != v {
					t.Errorf("Expected header key \"%s\" = \"%s\", but key does not exist or = \"%s\"", k, v, r)
				}
			}

			/* Response body check
			*/
			if tt.o.body != nil {
				responseBody, err := io.ReadAll(response.Body)
				require.NoError(t, err)
				defer response.Body.Close()
				
				if !bytes.Equal(responseBody, tt.o.body) {
					t.Errorf("Expected body \"%s\", got \"%s\"", tt.o.body, responseBody)
				}
			}

			/* URL store changes check
			*/
			assert.Equal(t, store.s, tt.o.URLStore)	
		})
	}	
}	

