package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/omekov/superapp-backend/pkg/grpc_errors"
)

func (r *userRepo) SetCacheUser(ctx context.Context, key string, seconds int, user *User) error {
	r.rw.Lock()
	defer r.rw.Unlock()

	if r.rdb == nil {
		return ErrNotConnection
	}

	userBytes, err := json.Marshal(user)
	if err != nil {
		return err
	}

	redisKey := fmt.Sprintf("%s:%s", userKey, key)
	return r.rdb.Set(ctx, redisKey, userBytes, time.Second*time.Duration(seconds)).Err()
}

func (r *userRepo) GetCacheByID(ctx context.Context, key string) (User, error) {
	r.rw.Lock()
	defer r.rw.Unlock()

	user := User{}
	if r.rdb == nil {
		return user, ErrNotConnection
	}

	redisKey := fmt.Sprintf("%s:%s", userKey, key)
	userBytes, err := r.rdb.Get(ctx, redisKey).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return user, grpc_errors.ErrNotFound
		}
		return user, err
	}

	if err = json.Unmarshal(userBytes, &user); err != nil {
		return user, err
	}

	return user, nil

}

func (r *userRepo) DeleteCacheUser(ctx context.Context, key string) error {
	r.rw.Lock()
	defer r.rw.Unlock()

	if r.rdb == nil {
		return ErrNotConnection
	}
	return r.rdb.Del(ctx, fmt.Sprintf("%s:%s", userKey, key)).Err()
}

func (r *userRepo) CreateSession(ctx context.Context, sess *Session, expire int) (string, error) {
	r.rw.Lock()
	defer r.rw.Unlock()

	sess.SessionID = uuid.New().String()
	if r.rdb == nil {
		return sess.SessionID, ErrNotConnection
	}

	sessionKey := fmt.Sprintf("%s:%s", sessionKey, sess.SessionID)

	sessBytes, err := json.Marshal(&sess)
	if err != nil {
		return "", err
	}

	if err = r.rdb.Set(ctx, sessionKey, sessBytes, time.Second*time.Duration(expire)).Err(); err != nil {
		return "", err
	}

	return sess.SessionID, nil
}

// Get session by id
func (r *userRepo) GetSessionByID(ctx context.Context, sessionID string) (Session, error) {
	r.rw.Lock()
	defer r.rw.Unlock()

	sess := Session{}
	if r.rdb == nil {
		return sess, ErrNotConnection
	}

	sessBytes, err := r.rdb.Get(ctx, fmt.Sprintf("%s:%s", sessionKey, sessionID)).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return sess, grpc_errors.ErrNotFound
		}
		return sess, err
	}

	if err = json.Unmarshal(sessBytes, &sess); err != nil {
		return sess, err
	}
	return sess, nil
}

func (r *userRepo) DeleteSessionByID(ctx context.Context, sessionID string) error {
	r.rw.Lock()
	defer r.rw.Unlock()

	if r.rdb == nil {
		return ErrNotConnection
	}

	err := r.rdb.Del(ctx, fmt.Sprintf("%s:%s", sessionKey, sessionID)).Err()
	return err
}
