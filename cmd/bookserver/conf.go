package main

import (
	"os"
)

// env vars
const (
	EnvPort = "APP_PORT"
	EnvDSN  = "DB_DSN"
)

// Config contains the required details to start the app
type Config struct {
	Port string
	DSN  string
}

// NewConfigFromEnv builds a Config from the env vars
func NewConfigFromEnv() Config {
	return Config{
		Port: os.Getenv(EnvPort),
		DSN:  os.Getenv(EnvDSN),
	}
}

// IsValid tells if the Config object is in a valid state
func (c Config) IsValid() bool {
	if c.Port == "" || c.DSN == "" {
		return false
	}
	return true
}
