package main

import (
	"container/list"
	"sync"

	"go.uber.org/zap"
)

// Pokemon represents the structure for a Pokemon
type Pokemon struct {
	ID        int      `json:"id"`
	Name      string   `json:"name"`
	Type      string   `json:"type"`
	Height    int      `json:"height"`
	Weight    int      `json:"weight"`
	Abilities []string `json:"abilities"`
}

// cacheEntry is a wrapper for cache entries
type cacheEntry struct {
	key   string
	value *Pokemon
}

// PokemonCache represents a LRU cache for Pokemon data
type PokemonCache struct {
	maxCapacity int
	cacheMap    map[string]*list.Element
	orderList   *list.List
	mu          sync.RWMutex
	logger      *zap.Logger
}

// NewPokemonCache creates a new PokemonCache with the given capacity
func NewPokemonCache(maxCapacity int, logger *zap.Logger) *PokemonCache {
	return &PokemonCache{
		maxCapacity: maxCapacity,
		cacheMap:    make(map[string]*list.Element),
		orderList:   list.New(),
		logger:      logger,
	}
}

// Set adds a new Pokemon to the cache or updates an existing one
func (c *PokemonCache) Set(key string, value *Pokemon) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.cacheMap[key]; ok {
		c.orderList.MoveToFront(elem)
		elem.Value.(*cacheEntry).value = value
		c.logger.Info("Updated Pokemon in cache", zap.String("key", key))
		return
	}

	if c.orderList.Len() >= c.maxCapacity {
		c.evict()
	}

	elem := c.orderList.PushFront(&cacheEntry{key, value})
	c.cacheMap[key] = elem
	c.logger.Info("Added Pokemon to cache", zap.String("key", key))
}

// Get retrieves a Pokemon from the cache by its name
func (c *PokemonCache) Get(key string) (*Pokemon, bool) {
	c.mu.Lock() //RLock()
	defer c.mu.Unlock()

	if elem, ok := c.cacheMap[key]; ok {
		c.orderList.MoveToFront(elem)
		c.logger.Info("Retrieved Pokemon from cache", zap.String("key", key))
		return elem.Value.(*cacheEntry).value, true
	}

	c.logger.Warn("Pokemon not found in cache", zap.String("key", key))
	return nil, false
}

// Delete removes a Pokemon from the cache by its name
func (c *PokemonCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.cacheMap[key]; ok {
		c.orderList.Remove(elem)
		delete(c.cacheMap, key)
		c.logger.Info("Deleted Pokemon from cache", zap.String("key", key))
	}
}

// evict removes the least recently used (LRU) item from the cache
func (c *PokemonCache) evict() {
	elem := c.orderList.Back()
	if elem != nil {
		c.orderList.Remove(elem)
		evictedKey := elem.Value.(*cacheEntry).key
		delete(c.cacheMap, evictedKey)
		c.logger.Info("Evicted Pokemon from cache", zap.String("key", evictedKey))
	}
}
