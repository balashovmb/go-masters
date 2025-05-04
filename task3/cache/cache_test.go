package cache

import (
	"testing"
)

func TestCache(t *testing.T) {
	cache := New[string, int]()
	cache.Set("key", 1000)
	if value, ok := cache.Get("key"); !ok || value != 1000 {
		t.Errorf("Got %d, want 1000", value)
	}
	if _, ok := cache.Get("key2"); ok {
		t.Errorf("Got %v, want false", ok)
	}
}
