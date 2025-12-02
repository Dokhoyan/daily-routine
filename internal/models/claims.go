package models

import (
	"github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	jwt.RegisteredClaims
	UserID string `json:"user_id"`
	Type   string `json:"type"` // "access" or "refresh"
}
