package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddGetFile(t *testing.T) {
	store, _ := NewFile("./test.txt")
	defer store.Close()
	var url string
	var err error

	// Try to get without adding first
	//
	_, err = store.Get("0")
	assert.Error(t, err)

	// Try to add few, get one back and check
	//
	i1, err := store.Add("http://www.google.com")
	assert.NoError(t, err)
	i2, err := store.Add("http://www.yandex.ru")
	assert.NoError(t, err)
	i3, err := store.Add("http://www.mail.ru")
	assert.NoError(t, err)
	url, err = store.Get(i1)
	assert.NoError(t, err)
	assert.Equal(t, url, "http://www.google.com")
	url, err = store.Get(i2)
	assert.NoError(t, err)
	assert.Equal(t, url, "http://www.yandex.ru")

	// Try to add a duplicate and check
	//
	id, err := store.Add("http://www.mail.ru")
	assert.NoError(t, err)
	assert.Equal(t, id, i3)
}
