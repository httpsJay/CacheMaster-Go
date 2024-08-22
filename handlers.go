package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

var pokeCache *PokemonCache
var logger *zap.Logger

func init() {
	// Initialize the logger
	logger, _ = zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any

	// Initialize the Pokemon cache with a max capacity of 100
	pokeCache = NewPokemonCache(100, logger)
	logger.Info("Initialized Pokemon cache with max capacity", zap.Int("capacity", 100))
}

// getPokemonByID handles GET requests to fetch a Pokemon by its ID
func getPokemonByID(w http.ResponseWriter, r *http.Request) {
	// Extract the ID from the URL
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid Pokemon ID", http.StatusBadRequest)
		logger.Error("Invalid Pokemon ID", zap.String("id", params["id"]))
		return
	}

	var found *Pokemon
	var foundKey string

	// Search for the Pokemon in the cache
	pokeCache.mu.Lock() // better to just RLock()
	for k, v := range pokeCache.cacheMap {
		if v.Value.(*cacheEntry).value.ID == id {
			found = v.Value.(*cacheEntry).value
			foundKey = k
			break
		}
	}
	pokeCache.mu.Unlock()

	if found == nil {
		http.Error(w, "Pokemon not found", http.StatusNotFound)
		logger.Warn("Pokemon not found with ID", zap.Int("id", id))
		return
	}

	// Move the found item to the front (most recently used)
	pokeCache.mu.Lock()
	pokeCache.orderList.MoveToFront(pokeCache.cacheMap[foundKey])
	pokeCache.mu.Unlock()

	// Encode the Pokemon data as JSON and send it in the response
	json.NewEncoder(w).Encode(found)
	logger.Info("Retrieved Pokemon by ID", zap.Int("id", id))
}

// getPokemonByName handles GET requests to fetch a Pokemon by its name
func getPokemonByName(w http.ResponseWriter, r *http.Request) {
	// Extract the name from the URL
	params := mux.Vars(r)
	name := params["name"]

	if name == "" {
		http.Error(w, "Pokemon name cannot be empty", http.StatusBadRequest)
		logger.Error("Pokemon name cannot be empty")
		return
	}

	// Search for the Pokemon in the cache
	pokemon, found := pokeCache.Get(name) //correct approach
	if !found {
		http.Error(w, "Pokemon not found", http.StatusNotFound)
		logger.Warn("Pokemon not found with Name", zap.String("name", name))
		return
	}

	// Encode the Pokemon data as JSON and send it in the response
	json.NewEncoder(w).Encode(pokemon)
	logger.Info("Retrieved Pokemon by Name", zap.String("name", name))
}

// deletePokemonByID handles DELETE requests to remove a Pokemon by its ID
func deletePokemonByID(w http.ResponseWriter, r *http.Request) {
	// Extract the ID from the URL
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid Pokemon ID", http.StatusBadRequest)
		logger.Error("Invalid Pokemon ID", zap.String("id", params["id"]))
		return
	}

	var foundKey string

	// Search for the Pokemon in the cache
	pokeCache.mu.Lock() //RLock()
	for k, v := range pokeCache.cacheMap {
		if v.Value.(*cacheEntry).value.ID == id {
			foundKey = k
			break
		}
	}
	pokeCache.mu.Unlock()

	if foundKey == "" {
		http.Error(w, "Pokemon not found", http.StatusNotFound)
		logger.Warn("Pokemon not found with ID", zap.Int("id", id))
		return
	}

	// Delete the Pokemon from the cache
	pokeCache.Delete(foundKey)
	w.WriteHeader(http.StatusOK)

}

// addPokemon handles POST requests to add a new Pokemon to the cache
func addPokemon(w http.ResponseWriter, r *http.Request) {
	var pokemon Pokemon

	// Decode the request body into a Pokemon struct
	err := json.NewDecoder(r.Body).Decode(&pokemon)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		logger.Error("Invalid request payload", zap.Error(err))
		return
	}

	// Validate the Pokemon data
	if pokemon.ID == 0 || pokemon.Name == "" || pokemon.Type == "" || pokemon.Height == 0 || pokemon.Weight == 0 || len(pokemon.Abilities) == 0 {
		http.Error(w, "Missing Pokemon data", http.StatusBadRequest)
		logger.Error("Missing Pokemon data", zap.Any("pokemon", pokemon))
		return
	}

	// Add the Pokemon to the cache
	pokeCache.Set(pokemon.Name, &pokemon)
	w.WriteHeader(http.StatusCreated)
}
