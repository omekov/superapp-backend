package domain

import (
	"time"

	"github.com/google/uuid"
)

// User ...
type User struct {
	ID       uuid.UUID
	Username string
	Password string
	Email    string
}

// Token ...
type Token struct {
	AccessToken  string
	RefreshToken string
}

// SessionLog ...
type SessionLog struct {
	ID           uuid.UUID `json:"id"`
	SessionID    uuid.UUID `json:"sessionID"`
	Username     string    `json:"username"`
	RefreshToken string    `json:"refreshToken"`
	UserAgent    string    `json:"userAgent"`
	ClientIP     string    `json:"client_ip"`
	IsActive     bool      `json:"isActive"`
	ExpiresAt    time.Time `json:"expiresAt"`
}
