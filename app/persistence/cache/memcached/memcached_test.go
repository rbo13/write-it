package memcached_test

import (
	"testing"

	"github.com/rbo13/write-it/app/persistence/cache"
	"github.com/rbo13/write-it/app/persistence/cache/memcached"
)

var testKey = "myTestKey1"

func TestSetCache(t *testing.T) {
	mem := memcached.New("localhost", "11211", "localhost:11211")

	ok, err := cache.Set(mem, testKey, "chardy")

	if err != nil {
		t.Error(err)
	}

	if !ok {
		t.Error("Error on saving")
	}
}
