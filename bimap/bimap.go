package bimap

import (
	"errors"
	"fmt"
)

// BiMap is a bidirectional map with unique keys and values
type BiMap[K comparable, V comparable] struct {
	keyToValue map[K]V
	valueToKey map[V]K
}

// NewBiMap initializes a new BiMap
func NewBiMap[K comparable, V comparable]() *BiMap[K, V] {
	return &BiMap[K, V]{
		keyToValue: make(map[K]V),
		valueToKey: make(map[V]K),
	}
}

// NewBiMapFromMap creates a BiMap from a standard map with duplicate checks.
func NewBiMapFromMap[K comparable, V comparable](input map[K]V) (*BiMap[K, V], error) {
	bimap := NewBiMap[K, V]()

	for key, value := range input {
		// Check for duplicate values
		if _, exists := bimap.valueToKey[value]; exists {
			return nil, fmt.Errorf("duplicate value found: %v", value)
		}
		// Add the key-value pair to the BiMap
		bimap.keyToValue[key] = value
		bimap.valueToKey[value] = key
	}

	return bimap, nil
}

// Set adds a key-value pair to the map
func (b *BiMap[K, V]) Set(key K, value V) error {
	// Check for existing key or value
	if _, exists := b.keyToValue[key]; exists {
		return errors.New("key already exists")
	}
	if _, exists := b.valueToKey[value]; exists {
		return errors.New("value already exists")
	}

	b.keyToValue[key] = value
	b.valueToKey[value] = key
	return nil
}

// GetByKey retrieves a value by its key
func (b *BiMap[K, V]) GetByKey(key K) (V, bool) {
	value, exists := b.keyToValue[key]
	return value, exists
}

// GetByValue retrieves a key by its value
func (b *BiMap[K, V]) GetByValue(value V) (K, bool) {
	key, exists := b.valueToKey[value]
	return key, exists
}

// DeleteByKey removes a key-value pair by its key
func (b *BiMap[K, V]) DeleteByKey(key K) bool {
	value, exists := b.keyToValue[key]
	if !exists {
		return false
	}

	delete(b.keyToValue, key)
	delete(b.valueToKey, value)
	return true
}

// DeleteByValue removes a key-value pair by its value
func (b *BiMap[K, V]) DeleteByValue(value V) bool {
	key, exists := b.valueToKey[value]
	if !exists {
		return false
	}

	delete(b.valueToKey, value)
	delete(b.keyToValue, key)
	return true
}
