package handlers

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nickeroshenkov/urlShortener/internal/app/storage"
)

type inputProvided struct {
	method   string
	url      string
	body     io.Reader
	URLStore []string
}

type outputDesired struct {
	code   int
	header map[string]string
	body   []byte
}

var tests = []struct {
	name string
	i    inputProvided
	o    outputDesired
}{
	{
		name: "Try to get with no arguments",
		i: inputProvided{
			method:   http.MethodGet,
			url:      "/",
			body:     nil,
			URLStore: nil,
		},
		o: outputDesired{
			code:   http.StatusBadRequest,
			header: nil,
			body:   nil,
		},
	},
	{
		name: "Try to get with no keys",
		i: inputProvided{
			method:   http.MethodGet,
			url:      "/",
			body:     nil,
			URLStore: nil,
		},
		o: outputDesired{
			code:   http.StatusBadRequest,
			header: nil,
			body:   nil,
		},
	},
	{
		name: "Try to get with a different key",
		i: inputProvided{
			method:   http.MethodGet,
			url:      "/?a=a",
			body:     nil,
			URLStore: nil,
		},
		o: outputDesired{
			code:   http.StatusBadRequest,
			header: nil,
			body:   nil,
		},
	},
	{
		name: "Try to get with a wrong value",
		i: inputProvided{
			method:   http.MethodGet,
			url:      "/?id=a",
			body:     nil,
			URLStore: nil,
		},
		o: outputDesired{
			code:   http.StatusBadRequest,
			header: nil,
			body:   nil,
		},
	},
	{
		name: "Try to get a non-existing URL #1",
		i: inputProvided{
			method:   http.MethodGet,
			url:      "/?id=0",
			body:     nil,
			URLStore: nil,
		},
		o: outputDesired{
			code:   http.StatusBadRequest,
			header: nil,
			body:   nil,
		},
	},
	{
		name: "Try to get a non-existing URL #2",
		i: inputProvided{
			method:   http.MethodGet,
			url:      "/?id=1",
			body:     nil,
			URLStore: []string{"http://www.google.com"},
		},
		o: outputDesired{
			code:   http.StatusBadRequest,
			header: nil,
			body:   nil,
		},
	},
	{
		name: "Get an existing full URL #1",
		i: inputProvided{
			method:   http.MethodGet,
			url:      "/?id=0",
			body:     nil,
			URLStore: []string{"http://www.google.com"},
		},
		o: outputDesired{
			code:   http.StatusTemporaryRedirect,
			header: map[string]string{"Location": "http://www.google.com"},
			body:   nil,
		},
	},
	{
		name: "Get an existing full URL #2",
		i: inputProvided{
			method:   http.MethodGet,
			url:      "/?id=1",
			body:     nil,
			URLStore: []string{"http://www.google.com", "http://www.yandex.ru"},
		},
		o: outputDesired{
			code:   http.StatusTemporaryRedirect,
			header: map[string]string{"Location": "http://www.yandex.ru"},
			body:   nil,
		},
	},
}

func TestUserViewHandler(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage.URLStore = tt.i.URLStore
			request := httptest.NewRequest(tt.i.method, tt.i.url, tt.i.body)
			response := httptest.NewRecorder()
			h := http.HandlerFunc(Shortener)
			h.ServeHTTP(response, request)
			result := response.Result()

			if result.StatusCode != tt.o.code {
				t.Errorf("Expected status code %d, but got %d", tt.o.code, response.Code)
			}

			for k, v := range tt.o.header {
				if r := result.Header.Get(k); r != v {
					t.Errorf("Expected header key \"%s\" = \"%s\", but key does not exist or = \"%s\"", k, v, r)
				}
			}

			if tt.o.body != nil {
				defer result.Body.Close()
				resultBody, err := io.ReadAll(result.Body)
				if err != nil {
					t.Fatal(err)
				}
				if !bytes.Equal(resultBody, tt.o.body) {
					t.Errorf("Expected body \"%s\", got \"%s\"", tt.o.body, resultBody)
				}
			}
		})
	}
}
