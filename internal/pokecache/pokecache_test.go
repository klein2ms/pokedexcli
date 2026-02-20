package pokecache

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestReapLoop(t *testing.T) {
	const baseTime = 5 * time.Millisecond
	const waitTime = baseTime + 5*time.Millisecond

	sut := NewCache(baseTime)
	err := sut.Add("https://example.com", []byte("testdata"))
	assert.NoError(t, err)

	_, ok := sut.Get("https://example.com")
	assert.True(t, ok)

	time.Sleep(waitTime)

	_, ok = sut.Get("https://example.com")
	assert.False(t, ok)
}

func TestAddGet(t *testing.T) {
	const interval = 5 * time.Second
	cases := []struct {
		key string
		val []byte
	}{
		{
			key: "https://example.com",
			val: []byte("testdata"),
		},
		{
			key: "https://example.com/path",
			val: []byte("moretestdata"),
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %v", i), func(t *testing.T) {
			cache := NewCache(interval)
			err := cache.Add(c.key, c.val)
			assert.NoError(t, err)

			val, ok := cache.Get(c.key)
			assert.True(t, ok)
			assert.Equal(t, c.val, val)
		})
	}
}
