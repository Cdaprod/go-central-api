package registry

import (
	"fmt"
	"plugin"
	"sync"

	"github.com/Cdaprod/go-central-api/config"
	"github.com/Cdaprod/go-central-api/integrations/repocate"
	"github.com/Cdaprod/go-central-api/integrations/minio"
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

func LoadService(r *APIRegistry, svc config.ServiceConfig) error {
	switch svc.Type {
	case "plugin":
		return loadPluginService(r, svc)
	case "builtin":
		return loadBuiltinService(r, svc)
	default:
		return fmt.Errorf("unknown service type: %s", svc.Type)
	}
}

func (r *APIRegistry) loadBuiltinService(svc config.ServiceConfig) error {
	var api API
	switch svc.Name {
	case "repocate":
		api = repocate.NewRepocateAPI(svc.URL)
	case "minio":
		api = minio.NewMinioAPI(svc.URL)
	default:
		return fmt.Errorf("unknown builtin service: %s", svc.Name)
	}
	r.Register(svc.Name, api)
	return nil
}

func loadPluginService(r *APIRegistry, svc config.ServiceConfig) error {
	pluginPath := svc.GetPluginPath()
	if pluginPath == "" {
		return fmt.Errorf("plugin path not specified for service %s", svc.Name)
	}

	p, err := plugin.Open(pluginPath)
	if err != nil {
		return fmt.Errorf("failed to open plugin for service %s: %v", svc.Name, err)
	}

	symNewAPI, err := p.Lookup("NewAPI")
	if err != nil {
		return fmt.Errorf("failed to find NewAPI symbol in plugin for service %s: %v", svc.Name, err)
	}

	newAPIFunc, ok := symNewAPI.(func(string) (API, error))
	if !ok {
		return fmt.Errorf("invalid NewAPI function signature in plugin for service %s", svc.Name)
	}

	api, err := newAPIFunc(svc.URL)
	if err != nil {
		return fmt.Errorf("failed to create API instance for service %s: %v", svc.Name, err)
	}

	r.Register(svc.Name, api)
	return nil
}

