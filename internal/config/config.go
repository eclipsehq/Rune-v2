package config

import "sync"

type Config struct {
	Token   string `json:"token"`
	OwnerID string `json:"owner_id"`
	Prefix  string `json:"prefix"`
}

var (
	Mu  sync.Mutex
	Cfg Config
)