package config

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type Config struct {
	Token   string `json:"token"`
	OwnerID string `json:"owner_id"`
	Prefix  string `json:"prefix"`
}

var (
	Mu  sync.Mutex
	Cfg Config
)

// LoadConfig reads the configuration file into the global Cfg variable.
func LoadConfig(path string) error {
	Mu.Lock()
	defer Mu.Unlock()

	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("could not open config: %w", err)
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(&Cfg); err != nil {
		return fmt.Errorf("could not decode config: %w", err)
	}

	return nil
}