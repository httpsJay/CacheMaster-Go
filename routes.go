package main

import (
    "github.com/gorilla/mux"
)

// setupRouter initializes the router and sets up the route handlers
func setupRouter() *mux.Router {
    router := mux.NewRouter()

    router.HandleFunc("/pokemon/id/{id}", getPokemonByID).Methods("GET")
    router.HandleFunc("/pokemon/name/{name}", getPokemonByName).Methods("GET")
    router.HandleFunc("/pokemon/id/{id}", deletePokemonByID).Methods("DELETE")
    router.HandleFunc("/pokemon", addPokemon).Methods("POST")

    return router
}
