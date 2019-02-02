package cache

// Cacher sets the basic caching functionality.
// E.g: Set, Get, Delete
type Cacher interface {
	Set(string, interface{}) (bool, error)
	Get(string) (interface{}, error)
	Delete(string) (bool, error)
}

// Set sets data to the cache
// and returns a boolean value
// data is saved successfully,
// returns error otherwise
func Set(c Cacher, key string, data interface{}) (bool, error) {
	return c.Set(key, data)
}

// Get retrieves data from cache
func Get(c Cacher, key string) (interface{}, error) {
	return c.Get(key)
}

// Delete deletes an item from the cache
// using the specified key. Returns
// boolean value if succcessful, error otherwise
func Delete(c Cacher, key string) (bool, error) {
	return c.Delete(key)
}
