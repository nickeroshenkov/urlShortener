package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddGet(t *testing.T) {
	store := NewURLStore()
	var url string
	var err error
	
	/* Try to get without adding first
	*/
	_, err = store.Get("0")
	assert.Error(t, err)

    /* Try to add few, get one back and check
	*/
	i1 := store.Add ("http://www.google.com")
	i2 := store.Add ("http://www.yandex.ru")
	i3 := store.Add ("http://www.mail.ru")
	url, err = store.Get(i1)
	assert.NoError(t, err)
	assert.Equal(t, url, "http://www.google.com")
	url, err = store.Get(i2)
	assert.NoError(t, err)
	assert.Equal(t, url, "http://www.yandex.ru")

    /* Try to add a duplicate and check
	*/
	id := store.Add ("http://www.mail.ru")
	assert.Equal(t, id, i3)
}
