package memcached_test

import (
	"testing"

	"github.com/rbo13/write-it/app/persistence/cache"
	"github.com/rbo13/write-it/app/persistence/cache/memcached"
)

var (
	testKey = "myTestKey1"
)

func TestMemcached(t *testing.T) {
	mem := memcached.New("localhost", "11211", "localhost:11211")

	t.Run("Set Cache", func(t *testing.T) {
		ok, err := cache.Set(mem, testKey, "chardy")

		if err != nil {
			t.Error(err)
		}

		if !ok {
			t.Error("Error on saving")
		}
	})

	t.Run("Get Cache", func(t *testing.T) {
		val, err := cache.Get(mem, testKey)

		if err != nil {
			t.Error(err)
		}

		t.Log(val)
	})

	t.Run("Delete Cache", func(t *testing.T) {
		ok, err := cache.Delete(mem, testKey)

		if err != nil {
			t.Log(err)
		}

		if !ok {
			t.Log("Error deleting value inside cache")
		}

		t.Log(ok)
	})

}
