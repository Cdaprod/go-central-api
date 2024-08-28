package config

import (
	"encoding/json"
	"os"
)

type ServiceConfig struct {
	Name    string            `json:"name"`
	Type    string            `json:"type"`
	URL     string            `json:"url"`
	Options map[string]string `json:"options"`
}

type Config struct {
	ServerAddress string          `json:"server_address"`
	DatabaseURL   string          `json:"database_url"`
	JWTSecret     string          `json:"jwt_secret"`
	Services      []ServiceConfig `json:"services"`
}

func Load() (*Config, error) {
	file, err := os.Open("config.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cfg Config
	if err := json.NewDecoder(file).Decode(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}