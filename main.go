package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"

	"github.com/Cdaprod/go-central-api/config"
	"github.com/Cdaprod/go-central-api/handlers"
	"github.com/Cdaprod/go-central-api/middleware"
	"github.com/Cdaprod/go-central-api/registry"
	"github.com/Cdaprod/go-central-api/utils"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize API registry
	apiRegistry := registry.NewAPIRegistry()

	// Load all services dynamically
	if err := loadServices(apiRegistry, cfg); err != nil {
		log.Fatalf("Failed to load services: %v", err)
	}
	log.Println("All services loaded successfully")

	// Create router
	r := mux.NewRouter()

	// Apply middleware
	r.Use(middleware.Logging)
	r.Use(middleware.Authentication(cfg))
	r.Use(middleware.CORS)

	// Create API Gateway
	gateway := handlers.NewAPIGateway(apiRegistry)

	// Register routes
	gateway.RegisterRoutes(r)

	// Add health check endpoint
	r.HandleFunc("/health", handlers.HealthCheckHandler).Methods("GET")

	// Add metrics endpoint
	r.HandleFunc("/metrics", handlers.MetricsHandler).Methods("GET")

	// Create server
	srv := &http.Server{
		Addr:         cfg.ServerAddress,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on %s", cfg.ServerAddress)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	utils.TimeTrack(time.Now(), "Server shutdown")
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}

func loadServices(apiRegistry *registry.APIRegistry, cfg *config.Config) error {
	for _, svc := range cfg.Services {
		if err := registry.LoadService(apiRegistry, svc); err != nil {
			return err
		}
		log.Printf("Loaded service: %s", svc.Name)
	}
	return nil
}
