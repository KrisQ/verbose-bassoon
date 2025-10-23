package pokecache

import (
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	cache := NewCache(5 * time.Millisecond)

	cache.Add("key1", []byte("value1"))

	val, ok := cache.Get("key1")
	if !ok || string(val) != "value1" {
		t.Error("expected to find key1 with value1")
	}

	time.Sleep(6 * time.Second)

	_, ok = cache.Get("key1")
	if ok {
		t.Error("expected key1 to be reaped")
	}
}
