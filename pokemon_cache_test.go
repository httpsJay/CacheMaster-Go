package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func setupLogger() *zap.Logger {
	logger, _ := zap.NewProduction()
	return logger
}

func TestSetAndGetPokemon(t *testing.T) {
	logger := setupLogger()
	defer logger.Sync()

	cache := NewPokemonCache(2, logger)

	p1 := &Pokemon{ID: 1, Name: "Bulbasaur"}
	p2 := &Pokemon{ID: 2, Name: "Ivysaur"}

	cache.Set(p1.Name, p1)
	cache.Set(p2.Name, p2)

	assert.Equal(t, 2, len(cache.cacheMap), "expected cache size 2")

	p, found := cache.Get(p1.Name)
	assert.True(t, found, "expected Bulbasaur to be in cache")
	assert.Equal(t, p1.ID, p.ID, "expected Bulbasaur ID to match")

	p, found = cache.Get(p2.Name)
	assert.True(t, found, "expected Ivysaur to be in cache")
	assert.Equal(t, p2.ID, p.ID, "expected Ivysaur ID to match")
}

func TestEvictPokemon(t *testing.T) {
	logger := setupLogger()
	defer logger.Sync()

	cache := NewPokemonCache(2, logger)

	p1 := &Pokemon{ID: 1, Name: "Bulbasaur"}
	p2 := &Pokemon{ID: 2, Name: "Ivysaur"}
	p3 := &Pokemon{ID: 3, Name: "Venusaur"}

	cache.Set(p1.Name, p1)
	cache.Set(p2.Name, p2)
	cache.Set(p3.Name, p3)

	assert.Equal(t, 2, len(cache.cacheMap), "expected cache size 2")

	_, found := cache.Get(p1.Name)
	assert.False(t, found, "expected Bulbasaur to be evicted")

	p, found := cache.Get(p2.Name)
	assert.True(t, found, "expected Ivysaur to be in cache")
	assert.Equal(t, p2.ID, p.ID, "expected Ivysaur ID to match")

	p, found = cache.Get(p3.Name)
	assert.True(t, found, "expected Venusaur to be in cache")
	assert.Equal(t, p3.ID, p.ID, "expected Venusaur ID to match")
}

func TestDeletePokemon(t *testing.T) {
	logger := setupLogger()
	defer logger.Sync()

	cache := NewPokemonCache(2, logger)

	p1 := &Pokemon{ID: 1, Name: "Bulbasaur"}
	p2 := &Pokemon{ID: 2, Name: "Ivysaur"}

	cache.Set(p1.Name, p1)
	cache.Set(p2.Name, p2)

	cache.Delete(p1.Name)
	_, found := cache.Get(p1.Name)
	assert.False(t, found, "expected Bulbasaur to be deleted")

	cache.Delete(p2.Name)
	_, found = cache.Get(p2.Name)
	assert.False(t, found, "expected Ivysaur to be deleted")
}

func TestCacheEvictionOrder(t *testing.T) {
	logger := setupLogger()
	defer logger.Sync()

	cache := NewPokemonCache(2, logger)

	p1 := &Pokemon{ID: 1, Name: "Bulbasaur"}
	p2 := &Pokemon{ID: 2, Name: "Ivysaur"}
	p3 := &Pokemon{ID: 3, Name: "Venusaur"}

	cache.Set(p1.Name, p1)
	cache.Set(p2.Name, p2)
	cache.Get(p1.Name) // Access Bulbasaur to make it recently used
	cache.Set(p3.Name, p3)

	assert.Equal(t, 2, len(cache.cacheMap), "expected cache size 2")

	_, found := cache.Get(p2.Name)
	assert.False(t, found, "expected Ivysaur to be evicted")

	p, found := cache.Get(p1.Name)
	assert.True(t, found, "expected Bulbasaur to be in cache")
	assert.Equal(t, p1.ID, p.ID, "expected Bulbasaur ID to match")

	p, found = cache.Get(p3.Name)
	assert.True(t, found, "expected Venusaur to be in cache")
	assert.Equal(t, p3.ID, p.ID, "expected Venusaur ID to match")
}
