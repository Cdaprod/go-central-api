package registry

import (
	"fmt"
	"plugin"
	"reflect"
	"sync"

	"github.com/Cdaprod/go-central-api/config"
)

type API interface {
	GetName() string
	Handle(method, path string, payload []byte) ([]byte, error)
}

type APIRegistry struct {
	apis map[string]API
	mu   sync.RWMutex
}

func NewAPIRegistry() *APIRegistry {
	return &APIRegistry{
		apis: make(map[string]API),
	}
}

func (r *APIRegistry) Register(name string, api API) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.apis[name] = api
}

func (r *APIRegistry) Get(name string) (API, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	api, ok := r.apis[name]
	return api, ok
}

func (r *APIRegistry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var names []string
	for name := range r.apis {
		names = append(names, name)
	}
	return names
}

func (r *APIRegistry) LoadServices(cfg *config.Config) error {
	for _, svc := range cfg.Services {
		switch svc.Type {
		case "builtin":
			if err := r.loadBuiltinService(svc); err != nil {
				return err
			}
		case "plugin":
			if err := r.loadPluginService(svc); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unknown service type: %s", svc.Type)
		}
	}
	return nil
}

func (r *APIRegistry) loadBuiltinService(svc config.ServiceConfig) error {
	apiType, exists := builtinAPIs[svc.Name]
	if !exists {
		return fmt.Errorf("builtin service not found: %s", svc.Name)
	}

	apiInstance := reflect.New(apiType).Interface().(API)
	r.Register(svc.Name, apiInstance)
	return nil
}

func (r *APIRegistry) loadPluginService(svc config.ServiceConfig) error {
	p, err := plugin.Open(svc.Options["path"])
	if err != nil {
		return fmt.Errorf("failed to open plugin: %v", err)
	}

	symAPI, err := p.Lookup("API")
	if err != nil {
		return fmt.Errorf("failed to find API symbol in plugin: %v", err)
	}

	api, ok := symAPI.(API)
	if !ok {
		return fmt.Errorf("invalid API implementation in plugin")
	}

	r.Register(svc.Name, api)
	return nil
}

var builtinAPIs = map[string]reflect.Type{
	// Register built-in API types here
	"repocate": reflect.TypeOf((*RepocateAPI)(nil)).Elem(),
	"minio":    reflect.TypeOf((*MinioAPI)(nil)).Elem(),
}

// Placeholder types for built-in APIs
type RepocateAPI struct{}
type MinioAPI struct{}

func (a *RepocateAPI) GetName() string { return "repocate" }
func (a *MinioAPI) GetName() string    { return "minio" }

// Implement Handle method for each API type
func (a *RepocateAPI) Handle(method, path string, payload []byte) ([]byte, error) {
	// Implement Repocate-specific logic
	return nil, nil
}

func (a *MinioAPI) Handle(method, path string, payload []byte) ([]byte, error) {
	// Implement MinIO-specific logic
	return nil, nil
}