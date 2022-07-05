package repository

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/omekov/superapp-backend/pkg/logger"
)

type Repository struct {
	User Userer
}

func NewRepository(rdb *redis.Client, db *sqlx.DB, log logger.Logger) *Repository {
	return &Repository{
		User: newUserRepo(rdb, db, log),
	}
}

type Userer interface {
	Create(ctx context.Context, user User) (uuid.UUID, error)
	GetByUsername(ctx context.Context, username string) (User, error)
	GetByID(ctx context.Context, id string) (User, error)
	SetCacheUser(ctx context.Context, key string, seconds int, user *User) error
	GetCacheByID(ctx context.Context, key string) (User, error)
	DeleteCacheUser(ctx context.Context, key string) error
	CreateSession(ctx context.Context, sess *Session, expire int) (string, error)
	GetSessionByID(ctx context.Context, sessionID string) (Session, error)
	DeleteSessionByID(ctx context.Context, sessionID string) error
	UpdatePassword(ctx context.Context, userID uuid.UUID, newPassword string) error
	GetByEmail(ctx context.Context, email string) (User, error)
	UpdateState(ctx context.Context, userID uuid.UUID, state, pinCode string) error
}
