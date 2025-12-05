package auth

import (
	"time"

	"github.com/Dokhoyan/daily-routine/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

const (
	accessTokenTTL  = 30 * time.Second // Временно уменьшено для тестирования
	refreshTokenTTL = 7 * 24 * time.Hour
)

func (s *serv) generateAccessToken(userID string) (string, error) {
	claims := &models.UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTokenTTL)),
		},
		UserID: userID,
		Type:   "access",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

func (s *serv) generateRefreshToken(userID string) (string, error) {
	claims := &models.UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshTokenTTL)),
		},
		UserID: userID,
		Type:   "refresh",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}
