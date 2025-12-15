package models

import "time"

type RefreshToken struct {
	ID         int64      `json:"id"`
	UserID     int64      `json:"user_id"`
	Token      string     `json:"token"` // Хэш токена
	ExpiresAt  time.Time  `json:"expires_at"`
	CreatedAt  time.Time  `json:"created_at"`
	RevokedAt  *time.Time `json:"revoked_at,omitempty"`
	DeviceInfo *string    `json:"device_info,omitempty"`
	IPAddress  *string    `json:"ip_address,omitempty"`
}

type TokenBlacklistEntry struct {
	ID        int64     `json:"id"`
	TokenHash string    `json:"token_hash"`
	UserID    int64     `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
	RevokedAt time.Time `json:"revoked_at"`
	Reason    *string   `json:"reason,omitempty"`
}
