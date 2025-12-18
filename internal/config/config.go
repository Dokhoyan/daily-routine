package config

import (
	"os"

	"github.com/joho/godotenv"
)

func Load(path string) error {
	err := godotenv.Load(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	return nil
}

type HTTPConfig interface {
	Address() string
}

type PGConfig interface {
	DSN() string
}

type ServerConfig interface {
	GetPort() string
}

type TelegramConfig interface {
	GetBotToken() string
}

type JWTConfig interface {
	GetSecret() string
}

type CORSConfig interface {
	GetAllowedOrigin() string
}

type TestConfig interface {
	IsTestModeEnabled() bool
	GetTestUserID() int64
}

type AdminConfig interface {
	GetUsername() string
	GetPassword() string
}
