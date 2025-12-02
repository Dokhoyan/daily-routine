package config

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
)

const (
	telegramBotTokenEnvName = "TELEGRAM_BOT_TOKEN"
)

type telegramConfig struct {
	botToken string
}

func NewTelegramConfig() (TelegramConfig, error) {
	botToken := os.Getenv(telegramBotTokenEnvName)

	testCfg := NewTestConfig()
	if botToken == "" {
		if testCfg.IsTestModeEnabled() {
			botToken = "test_bot_token"
			fmt.Println("⚠️  WARNING: Using default test value for TELEGRAM_BOT_TOKEN")
			fmt.Println("⚠️  This is only for development/testing. Set proper value in .env for production!")
		} else {
			return nil, errors.New("telegram bot token not found")
		}
	}

	return &telegramConfig{
		botToken: botToken,
	}, nil
}

func (cfg *telegramConfig) GetBotToken() string {
	return cfg.botToken
}
