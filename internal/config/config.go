package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	DataDirectory string `json:"data_directory"`
	ListenAddress string `json:"listen_address"`
}

func Load(path string) (*Config, error) {
	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	var cfg Config

	if err = json.NewDecoder(file).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("decode config %q, %w", path, err)
	}

	return &cfg, nil
}
