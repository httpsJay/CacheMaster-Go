package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestAddPokemon(t *testing.T) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	pokeCache = NewPokemonCache(100, logger)

	router := setupRouter()

	validPokemon := Pokemon{ID: 1, Name: "Bulbasaur", Type: "Grass/Poison", Height: 7, Weight: 69, Abilities: []string{"Overgrow", "Chlorophyll"}}
	payload, _ := json.Marshal(validPokemon)

	req, _ := http.NewRequest("POST", "/pokemon", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code, "handler returned wrong status code")

	// Test invalid JSON
	invalidPayload := []byte(`{"id": 2, "name": "Ivysaur", "type": Grass/Poison"}`)
	req, _ = http.NewRequest("POST", "/pokemon", bytes.NewBuffer(invalidPayload))
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code, "handler returned wrong status code for invalid JSON")

	// Test missing data
	invalidPokemon := Pokemon{ID: 2, Name: "", Type: "Grass/Poison", Height: 7, Weight: 69, Abilities: []string{"Overgrow"}}
	payload, _ = json.Marshal(invalidPokemon)
	req, _ = http.NewRequest("POST", "/pokemon", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code, "handler returned wrong status code for missing data")
}

func TestGetPokemonByID(t *testing.T) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	pokeCache = NewPokemonCache(100, logger)

	pokemon := Pokemon{ID: 1, Name: "Bulbasaur", Type: "Grass/Poison", Height: 7, Weight: 69, Abilities: []string{"Overgrow", "Chlorophyll"}}
	pokeCache.Set(pokemon.Name, &pokemon)

	router := setupRouter()

	req, _ := http.NewRequest("GET", "/pokemon/id/1", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "handler returned wrong status code")

	var returnedPokemon Pokemon
	err := json.NewDecoder(rr.Body).Decode(&returnedPokemon)
	assert.NoError(t, err, "could not decode response")

	assert.Equal(t, pokemon.ID, returnedPokemon.ID, "handler returned unexpected Pokemon ID")

	// Test invalid ID
	req, _ = http.NewRequest("GET", "/pokemon/id/invalid", nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code, "handler returned wrong status code for invalid ID")

	// Test non-existent ID
	req, _ = http.NewRequest("GET", "/pokemon/id/999", nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code, "handler returned wrong status code for non-existent ID")
}

func TestGetPokemonByName(t *testing.T) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	pokeCache = NewPokemonCache(100, logger)

	pokemon := Pokemon{ID: 1, Name: "Bulbasaur", Type: "Grass/Poison", Height: 7, Weight: 69, Abilities: []string{"Overgrow", "Chlorophyll"}}
	pokeCache.Set(pokemon.Name, &pokemon)

	router := setupRouter()

	req, _ := http.NewRequest("GET", "/pokemon/name/Bulbasaur", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "handler returned wrong status code")

	var returnedPokemon Pokemon
	err := json.NewDecoder(rr.Body).Decode(&returnedPokemon)
	assert.NoError(t, err, "could not decode response")

	assert.Equal(t, pokemon.Name, returnedPokemon.Name, "handler returned unexpected Pokemon Name")

	// Test non-existent Name
	req, _ = http.NewRequest("GET", "/pokemon/name/NonExistent", nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code, "handler returned wrong status code for non-existent Name")

}

func TestDeletePokemonByID(t *testing.T) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	pokeCache = NewPokemonCache(100, logger)

	pokemon := Pokemon{ID: 1, Name: "Bulbasaur", Type: "Grass/Poison", Height: 7, Weight: 69, Abilities: []string{"Overgrow", "Chlorophyll"}}
	pokeCache.Set(pokemon.Name, &pokemon)

	router := setupRouter()

	req, _ := http.NewRequest("DELETE", "/pokemon/id/1", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "handler returned wrong status code")

	_, found := pokeCache.Get(pokemon.Name)
	assert.False(t, found, "expected Pokemon to be deleted")

	// Test invalid ID
	req, _ = http.NewRequest("DELETE", "/pokemon/id/invalid", nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code, "handler returned wrong status code for invalid ID")

	// Test non-existent ID
	req, _ = http.NewRequest("DELETE", "/pokemon/id/999", nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code, "handler returned wrong status code for non-existent ID")
}
