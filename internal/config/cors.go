package config

import (
	"os"
)

const (
	corsAllowedOriginEnvName = "ALLOWED_ORIGIN"
)

type corsConfig struct {
	allowedOrigin string
}

func NewCORSConfig() CORSConfig {
	return &corsConfig{
		allowedOrigin: os.Getenv(corsAllowedOriginEnvName),
	}
}

func (cfg *corsConfig) GetAllowedOrigin() string {
	return cfg.allowedOrigin
}
