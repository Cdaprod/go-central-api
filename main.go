package main

import (
	"log"
	"net/http"

	"github.com/Cdaprod/go-central-api/config"
	"github.com/Cdaprod/go-central-api/handlers"
	"github.com/Cdaprod/go-central-api/middleware"
	"github.com/Cdaprod/go-central-api/registry"
	"github.com/gorilla/mux"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize API registry
	apiRegistry := registry.NewAPIRegistry()

	// Load services dynamically
	if err := apiRegistry.LoadServices(cfg); err != nil {
		log.Fatalf("Failed to load services: %v", err)
	}

	// Create router
	r := mux.NewRouter()

	// Apply middleware
	r.Use(middleware.Logging)
	r.Use(middleware.Authentication)

	// Create API Gateway
	gateway := handlers.NewAPIGateway(apiRegistry)

	// Register routes
	gateway.RegisterRoutes(r)

	// Start server
	log.Printf("Starting server on %s", cfg.ServerAddress)
	if err := http.ListenAndServe(cfg.ServerAddress, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}