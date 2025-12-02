package config

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
)

const (
	jwtSecretEnvName = "JWT_SECRET"
)

type jwtConfig struct {
	secret string
}

func NewJWTConfig() (JWTConfig, error) {
	secret := os.Getenv(jwtSecretEnvName)

	testCfg := NewTestConfig()
	if secret == "" {
		if testCfg.IsTestModeEnabled() {
			secret = "test_jwt_secret_key_for_development_only"
			fmt.Println("⚠️  WARNING: Using default test value for JWT_SECRET")
			fmt.Println("⚠️  This is only for development/testing. Set proper value in .env for production!")
		} else {
			return nil, errors.New("jwt secret not found")
		}
	}

	return &jwtConfig{
		secret: secret,
	}, nil
}

func (cfg *jwtConfig) GetSecret() string {
	return cfg.secret
}
