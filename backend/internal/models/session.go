package models

import "time"

type Session struct {
	ID           int
	UserID       int
	RefreshToken string
	UserAgent    string
	IPAddress    string
	ExpiredAt    time.Time
	CreatedAt    time.Time
}

type SessionResponse struct {
	User
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
