package main

import (
	"net/http"

	"go.uber.org/zap"
)

func main() {
	// Setup the router
	router := setupRouter()
	logger.Info("Starting server on :8080")

	// Start the HTTP server
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		logger.Fatal("Server failed to start", zap.Error(err))
	}
}
