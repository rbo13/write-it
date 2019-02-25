package memcached_test

import (
	"encoding/json"
	"testing"

	"github.com/rbo13/write-it/app/persistence/cache"
	"github.com/rbo13/write-it/app/persistence/cache/memcached"
)

var (
	testKey = "myTestKey1"
)

type testData struct {
	Val string `json:"val"`
}

func TestMemcached(t *testing.T) {
	mem := memcached.New("localhost", "11211", "localhost:11211")

	t.Run("Set Cache", func(t *testing.T) {

		val := testData{
			Val: "Hello World",
		}

		d, err := json.Marshal(val)

		if err != nil {
			t.Error(err)
		}

		ok, err := cache.Set(mem, testKey, string(d))

		if err != nil {
			t.Error(err)
		}

		if !ok {
			t.Error("Error on saving")
		}
	})

	t.Run("Get Cache", func(t *testing.T) {
		var d testData

		err := cache.Get(mem, testKey, d)

		if err != nil {
			t.Error(err)
		}

		t.Log(d)
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
