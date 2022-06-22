package domain

import (
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
}
