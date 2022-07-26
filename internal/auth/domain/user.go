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
type UserSessionLog struct {
	SessionID   string
	Username    string
	UserAgent   string
	ClientIP    string
	HTTPMethod  string
	HTTPPath    string
	HTTPStatus  uint32
	HTTPReqBody string
	HTTPResBody string
	ExpiresAt   time.Time
}
