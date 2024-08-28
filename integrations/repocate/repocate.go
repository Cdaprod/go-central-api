package repocate

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type RepocateAPI struct {
	BaseURL string
}

func NewRepocateAPI(baseURL string) *RepocateAPI {
	return &RepocateAPI{BaseURL: baseURL}
}

func (a *RepocateAPI) GetName() string {
	return "repocate"
}

func (a *RepocateAPI) Handle(method, path string, payload []byte) ([]byte, error) {
	switch {
	case method == "GET" && path == "containers":
		return a.listContainers()
	// Add more cases for different endpoints
	default:
		return nil, fmt.Errorf("unknown endpoint: %s %s", method, path)
	}
}

func (a *RepocateAPI) listContainers() ([]byte, error) {
	resp, err := http.Get(a.BaseURL + "/containers")
	if err != nil {
		return nil, fmt.Errorf("failed to get containers from Repocate: %v", err)
	}
	defer resp.Body.Close()

	var containers []string
	if err := json.NewDecoder(resp.Body).Decode(&containers); err != nil {
		return nil, fmt.Errorf("failed to decode Repocate response: %v", err)
	}

	return json.Marshal(containers)
}
