package config

import (
	"os"
	"strconv"
)

const (
	maxSessionsEnvName = "MAX_ACTIVE_SESSIONS"
)

type authConfig struct {
	maxActiveSessions int
}

func NewAuthConfig() AuthConfig {
	maxSessions := 5 // default
	if maxSessionsStr := os.Getenv(maxSessionsEnvName); maxSessionsStr != "" {
		if parsed, err := strconv.Atoi(maxSessionsStr); err == nil && parsed > 0 {
			maxSessions = parsed
		}
	}

	return &authConfig{
		maxActiveSessions: maxSessions,
	}
}

func (cfg *authConfig) GetMaxActiveSessions() int {
	return cfg.maxActiveSessions
}

type AuthConfig interface {
	GetMaxActiveSessions() int
}
