package auth

import (
	"net/http"

	"github.com/Dokhoyan/daily-routine/internal/config"
	"github.com/Dokhoyan/daily-routine/internal/repository"
	"github.com/Dokhoyan/daily-routine/internal/service"
)

type serv struct {
	telegramBotToken string
	jwtSecret        string
	userRepo         repository.UserRepository
	tokenRepo        repository.TokenRepository
	authConfig       config.AuthConfig
	tokenCache       TokenCache
}

func NewService(telegramBotToken, jwtSecret string, userRepo repository.UserRepository, tokenRepo repository.TokenRepository, authConfig config.AuthConfig, tokenCache TokenCache) service.AuthService {
	return &serv{
		telegramBotToken: telegramBotToken,
		jwtSecret:        jwtSecret,
		userRepo:         userRepo,
		tokenRepo:        tokenRepo,
		authConfig:       authConfig,
		tokenCache:       tokenCache,
	}
}

func GetDeviceInfo(r *http.Request) string {
	userAgent := r.UserAgent()
	if userAgent == "" {
		return "unknown"
	}
	return userAgent
}

func GetIPAddress(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.Header.Get("X-Real-Ip")
	}
	if ip == "" {
		ip = r.RemoteAddr
	}
	return ip
}
