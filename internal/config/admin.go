package config

import (
	"errors"
	"os"
)

type adminConfig struct {
	username string
	password string
}

func NewAdminConfig() (AdminConfig, error) {
	username := os.Getenv("ADMIN_USERNAME")
	if username == "" {
		return nil, errors.New("ADMIN_USERNAME is not set")
	}

	password := os.Getenv("ADMIN_PASSWORD")
	if password == "" {
		return nil, errors.New("ADMIN_PASSWORD is not set")
	}

	return &adminConfig{
		username: username,
		password: password,
	}, nil
}

func (c *adminConfig) GetUsername() string {
	return c.username
}

func (c *adminConfig) GetPassword() string {
	return c.password
}

