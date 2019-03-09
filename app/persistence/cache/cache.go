package cache

import (
	"encoding/json"
	"errors"
)

var (
	errMissingKey = errors.New("err: missing key")
)

// Cacher sets the basic caching functionality.
// E.g. Set, Get, Delete
type Cacher interface {
	Set(string, string) (bool, error)
	Get(string) (string, error)
	Delete(string) (bool, error)
}

// Set sets data to the cache
// and returns a boolean value
// data is saved successfully,
// returns error otherwise
func Set(c Cacher, key string, data interface{}) (bool, error) {
	if key == "" {
		return false, errMissingKey
	}

	val, err := json.Marshal(data)

	if err != nil {
		return false, err
	}
	// _, err = cache.Set(mem, cacheKey, string(val))
	_, err = c.Set(key, string(val))

	if err != nil {
		return false, err
	}

	return true, nil
}

// Get retrieves data from cache and wraps the json.Unmarshal function.
func Get(c Cacher, key string, dest interface{}) error {
	if key == "" {
		return errMissingKey
	}

	data, err := c.Get(key)

	if err == nil && data != "" {
		return json.Unmarshal([]byte(data), &dest)
	}

	return err
}

// Delete deletes an item from the cache
// using the specified key. Returns
// boolean value if succcessful, error otherwise
func Delete(c Cacher, key string) (bool, error) {
	if key == "" {
		return false, errMissingKey
	}
	return c.Delete(key)
}

// Unmarshal wraps the json.Unmarshal to unmarshal the values inside the cache.
func Unmarshal(data string, dest interface{}) error {
	return json.Unmarshal([]byte(data), &dest)
}
