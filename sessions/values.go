package sessions

import "encoding/json"

// Values contains keys and associated values.
type Values interface {
	// Get gets the value associated with key. If there is no value associated
	// with key, Get returns an empty string.
	Get(key string) string

	// GetAll returns all keys and their associated value.
	GetAll() map[string]string

	// Remove removes values associated with the provided keys.
	Remove(keys ...string)

	// RemoveAll removes all keys and values.
	RemoveAll()

	// Set sets the key to value.
	Set(key, value string)

	// SetAll sets all provided keys to their associated value.
	SetAll(map[string]string)
}

// values is an unexported type that implements the Values interface.
type values map[string]string

// NewValues returns a new instance of Values.
func NewValues() Values {
	return make(values)
}

// Get gets the value associated with key. If there is no value associated with
// key, Get returns an empty string.
func (v values) Get(key string) string {
	if value, ok := v[key]; ok {
		return value
	}
	return ""
}

// GetAll returns all key-value pairs.
func (v values) GetAll() map[string]string {
	return map[string]string(v)
}

// Remove removes values associated with the keys.
func (v values) Remove(keys ...string) {
	for _, key := range keys {
		delete(v, key)
	}
}

// RemoveAll removes all keys and values.
func (v values) RemoveAll() {
	for key := range v {
		delete(v, key)
	}
}

// Set sets the key to value. It replaces an existing value.
func (v values) Set(key, value string) {
	v[key] = value
}

// SetAll sets all provided keys to their associated value.
func (v values) SetAll(pairs map[string]string) {
	for key, value := range pairs {
		v.Set(key, value)
	}
}

// ValuesFromJSON JSON decodes a map of key-value pairs. The result can be used
// as input for Values.SetAll.
func ValuesFromJSON(data []byte) (map[string]string, error) {
	temp := map[string]string{}
	if err := json.Unmarshal(data, &temp); err != nil {
		return nil, err
	}
	return temp, nil
}
