package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/Cdaprod/go-central-api/registry"
	"github.com/gorilla/mux"
)

type APIGateway struct {
	registry *registry.APIRegistry
}

func NewAPIGateway(registry *registry.APIRegistry) *APIGateway {
	return &APIGateway{registry: registry}
}

func (g *APIGateway) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/api/health", g.HealthCheckHandler).Methods("GET")
	r.HandleFunc("/api/{service}/{path:.*}", g.ProxyHandler).Methods("GET", "POST", "PUT", "DELETE")
}

func (g *APIGateway) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}

func (g *APIGateway) ProxyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	serviceName := vars["service"]
	path := vars["path"]

	api, ok := g.registry.Get(serviceName)
	if !ok {
		http.Error(w, "Service not found", http.StatusNotFound)
		return
	}

	// Read the request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	// Call the service's Handle method
	response, err := api.Handle(r.Method, path, body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write the response
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

// ListServicesHandler returns a list of all registered services
func (g *APIGateway) ListServicesHandler(w http.ResponseWriter, r *http.Request) {
	services := g.registry.List()
	json.NewEncoder(w).Encode(map[string][]string{"services": services})
}
