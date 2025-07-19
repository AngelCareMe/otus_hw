package hw04lrucache

import (
	"testing"

	"github.com/stretchr/testify/require" //nolint:depguard
)

func TestCahce_Clear(t *testing.T) {
	c := NewCache(3)
	c.Set("x", 10)
	c.Set("y", 20)
	c.Set("z", 30)

	require.Equal(t, 3, c.(*lruCache).queue.Len())

	c.Clear()

	require.Equal(t, 0, c.(*lruCache).queue.Len())
	require.Empty(t, c.(*lruCache).items)

	_, ok := c.Get("x")
	require.False(t, ok)
	_, ok = c.Get("y")
	require.False(t, ok)
	_, ok = c.Get("z")
	require.False(t, ok)

	ok = c.Set("new", 999)
	require.False(t, ok)
	val, ok := c.Get("new")
	require.True(t, ok)
	require.Equal(t, 999, val)
}
