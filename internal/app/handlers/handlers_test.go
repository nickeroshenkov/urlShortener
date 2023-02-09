package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
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
	URLStore []string
}

type outputDesired struct {
	code   int
	header map[string]string
	body   []byte
	URLStore []string
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
			URLStore: nil,
		},
		o: outputDesired{
			code:   http.StatusCreated,
			header: nil,
			body:   nil,
			URLStore: []string{"http://www.google.com"},
		},
	},
	{
		name: "Try to get with a wrong value",
		i: inputProvided{
			method:   http.MethodGet,
			path:     "/a",
			body:     nil,
			URLStore: nil,
		},
		o: outputDesired{
			code:   http.StatusBadRequest,
			header: nil,
			body:   nil,
			URLStore: nil,
		},
	},
	{
		name: "Try to get a non-existing URL #1",
		i: inputProvided{
			method:   http.MethodGet,
			path:     "/0",
			body:     nil,
			URLStore: nil,
		},
		o: outputDesired{
			code:   http.StatusBadRequest,
			header: nil,
			body:   nil,
			URLStore: nil,
		},
	},
	{
		name: "Try to get a non-existing URL #2",
		i: inputProvided{
			method:   http.MethodGet,
			path:     "/1",
			body:     nil,
			URLStore: []string{"http://www.google.com"},
		},
		o: outputDesired{
			code:   http.StatusBadRequest,
			header: nil,
			body:   nil,
			URLStore: []string{"http://www.google.com"},
		},
	},
	{
		name: "Get an existing full URL #1",
		i: inputProvided{
			method:   http.MethodGet,
			path:     "/0",
			body:     nil,
			URLStore: []string{"http://www.google.com"},
		},
		o: outputDesired{
			code:   http.StatusTemporaryRedirect,
			header: map[string]string{"Location": "http://www.google.com"},
			body:   nil,
			URLStore: []string{"http://www.google.com"},
		},
	},
	{
		name: "Get an existing full URL #2",
		i: inputProvided{
			method:   http.MethodGet,
			path:     "/1",
			body:     nil,
			URLStore: []string{"http://www.google.com", "http://www.yandex.ru"},
		},
		o: outputDesired{
			code:   http.StatusTemporaryRedirect,
			header: map[string]string{"Location": "http://www.yandex.ru"},
			body:   nil,
			URLStore: []string{"http://www.google.com", "http://www.yandex.ru"},
		},
	},
}

type urlStoreMock struct {
	urls []string
}

func (store *urlStoreMock) Add (url string) int {
	store.urls = append (store.urls, url)
	return len(store.urls)-1
}

func (store *urlStoreMock) Get (id int) (string, error) {
	if id >= len(store.urls) {
		return "", errors.New ("URL does not exist in the store")
	}
	return store.urls[id], nil
}

func TestSetRoute(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := urlStoreMock { urls: tt.i.URLStore }
			router := chi.NewRouter()
			SetRoute (&store, router)
			server := httptest.NewServer(router)
			defer server.Close()
			
			request, err := http.NewRequest(tt.i.method, server.URL+tt.i.path, bytes.NewReader(tt.i.body))
			require.NoError(t, err)
			// Provide CheckRedirect() to modify client's behavior for re-directs
			client := &http.Client {
				CheckRedirect: func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				},
			}
			response, err := client.Do(request)
			require.NoError(t, err)
			
			if response.StatusCode != tt.o.code {
				t.Errorf("Expected status code %d, but got %d", tt.o.code, response.StatusCode)
			}

			for k, v := range tt.o.header {
				if r := response.Header.Get(k); r != v {
					t.Errorf("Expected header key \"%s\" = \"%s\", but key does not exist or = \"%s\"", k, v, r)
				}
			}

			/* Response body check
			if tt.o.body != nil {
				responseBody, err := io.ReadAll(response.Body)
				require.NoError(t, err)
				defer response.Body.Close()
				
				if !bytes.Equal(responseBody, tt.o.body) {
					t.Errorf("Expected body \"%s\", got \"%s\"", tt.o.body, responseBody)
				}
			} */

			/* Check changes in URL store
			*/
			assert.Equal(t, store.urls, tt.o.URLStore)	
		})
	}	
}	
