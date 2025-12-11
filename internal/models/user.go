package models

import "time"

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	FirstName string    `json:"first_name"`
	PhotoURL  string    `json:"photo_url"`
	AuthDate  time.Time `json:"auth_date"`
	TokenTG   string    `json:"tokentg"`
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type AuthResponse struct {
	User      *User      `json:"user"`
	TokenPair *TokenPair `json:"tokens"`
}

type UserSettings struct {
	UserID       int64    `json:"user_id"`
	Timezone     string   `json:"timezone"`
	DoNotDisturb bool     `json:"do_not_disturb"`
	NotifyTimes  []string `json:"notify_times"`
}
