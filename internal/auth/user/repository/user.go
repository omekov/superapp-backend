package repository

import (
	"context"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/omekov/superapp-backend/pkg/logger"
	"github.com/pkg/errors"
)

const (
	sessionKey            = "sessions"
	userKey               = "user"
	userStateEnabled      = "enabled"
	userStateDisabled     = "disabled"
	userStateNotActivated = "notactivated"
)

var (
	ErrNotConnection = errors.New("no connection to database")
	ErrNoRowsUpdated = errors.New("failed to update because not found")
)

type userRepo struct {
	rdb *redis.Client
	db  *sqlx.DB
	log logger.Logger
	rw  sync.RWMutex
}

func newUserRepo(rdb *redis.Client, db *sqlx.DB, log logger.Logger) *userRepo {
	return &userRepo{rdb, db, log, sync.RWMutex{}}
}

func (r *userRepo) Create(ctx context.Context, user User) (uuid.UUID, error) {
	const queryCreate = `INSERT into users 
	(
		id,
		username, 
		password,
		email,
		state,
		pin_code
	) VALUES ($1,$2,$3,$4,$5,$6)
	RETURNING id`

	r.rw.Lock()
	defer r.rw.Unlock()

	var id uuid.UUID
	if r.db == nil {
		return id, ErrNotConnection
	}

	err := r.db.QueryRowContext(ctx, queryCreate, user.ID, user.UserName, user.Password, user.Email, user.State, user.PinCode).Scan(&id)
	if err != nil {
		return id, err
	}

	return id, nil
}

func (r *userRepo) GetByUsername(ctx context.Context, username string) (User, error) {
	const queryGetByUsername = `SELECT
	id,
	username,
	password,
	email,
	state
	FROM users
	WHERE username=$1
	ORDER by id DESC LIMIT 1`

	user := User{}
	if r.db == nil {
		return user, ErrNotConnection
	}

	err := r.db.GetContext(ctx, &user, queryGetByUsername, username)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (r *userRepo) GetByID(ctx context.Context, id string) (User, error) {
	const queryGetByID = `SELECT
	id,
	username,
	password,
	email,
	state
	FROM users
	WHERE id=$1
	ORDER by id DESC LIMIT 1`

	user := User{}
	if r.db == nil {
		return user, ErrNotConnection
	}

	err := r.db.GetContext(ctx, &user, queryGetByID, id)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *userRepo) UpdatePassword(ctx context.Context, userID uuid.UUID, newPassword string) error {
	const queryUpdatePassword = `UPDATE users
	SET password = $2
	WHERE id=$1`

	r.rw.Lock()
	defer r.rw.Unlock()

	if r.db == nil {
		return ErrNotConnection
	}

	result, err := r.db.ExecContext(ctx, queryUpdatePassword, userID, newPassword)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrNoRowsUpdated
	}

	return nil
}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (User, error) {
	const queryGetByEmail = `SELECT
	id,
	username,
	password,
	email,
	state,
	pin_code
	FROM users
	WHERE email=$1
	ORDER by id DESC LIMIT 1`

	user := User{}
	if r.db == nil {
		return user, ErrNotConnection
	}

	err := r.db.GetContext(ctx, &user, queryGetByEmail, email)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *userRepo) UpdateState(ctx context.Context, userID uuid.UUID, state, pinCode string) error {
	const queryUpdatePassword = `UPDATE users
	SET state = $2, pin_code = $3
	WHERE id=$1`

	r.rw.Lock()
	defer r.rw.Unlock()

	if r.db == nil {
		return ErrNotConnection
	}

	result, err := r.db.ExecContext(ctx, queryUpdatePassword, userID, state, pinCode)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrNoRowsUpdated
	}

	return nil
}
