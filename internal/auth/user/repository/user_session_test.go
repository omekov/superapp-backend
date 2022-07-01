package repository_test

import (
	"context"
	"testing"
	"time"

	. "github.com/omekov/superapp-backend/internal/auth/user/repository"
	"github.com/omekov/superapp-backend/pkg/grpc_errors"
	"github.com/omekov/superapp-backend/pkg/logger"
	"github.com/stretchr/testify/require"
)

func TestUser_SetCacheUser(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	type args struct {
		key     string
		seconds int
		user    User
	}
	user := TestGetUser(t)
	logging := logger.NewAPILogger("debug")
	testCases := []struct {
		name           string
		wantCreateUser bool
		wantErr        bool
		args           args
		repo           *Repository
		err            error
	}{
		{
			"valid",
			true,
			false,
			args{
				key:     user.ID.String(),
				seconds: 0,
				user:    user,
			},
			NewRepository(rdb, dbx, logging),
			nil,
		},
		{
			"invalid rdb is nil",
			false,
			true,
			args{
				key:     user.ID.String(),
				seconds: 0,
				user:    user,
			},
			NewRepository(nil, nil, logging),
			ErrNotConnection,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.repo.User.SetCacheUser(ctx, tc.args.key, tc.args.seconds, &tc.args.user)
			if err != nil {
				if tc.wantErr {
					require.Equal(t, err, tc.err)
				} else {
					require.Error(t, err)
				}
			}
		})
	}
}

func TestUser_GetCacheByID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	type args struct {
		id string
	}
	expUser := TestGetUser(t)
	logging := logger.NewAPILogger("debug")
	testCases := []struct {
		name           string
		wantCreateUser bool
		wantErr        bool
		args           args
		repo           *Repository
		err            error
	}{
		{
			"valid",
			true,
			false,
			args{
				id: expUser.ID.String(),
			},
			NewRepository(rdb, dbx, logging),
			nil,
		},
		{
			"invalid rdb is nil",
			false,
			true,
			args{
				id: expUser.ID.String(),
			},
			NewRepository(nil, nil, logging),
			ErrNotConnection,
		},
		{
			"invalid rdb userId not found",
			true,
			true,
			args{
				id: expUser.ID.String() + "test",
			},
			NewRepository(rdb, dbx, logging),
			grpc_errors.ErrNotFound,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			user, err := tc.repo.User.GetCacheByID(ctx, tc.args.id)
			if err != nil {
				if tc.wantErr {
					require.Equal(t, err, tc.err)
				} else {
					require.Error(t, err)
				}
			} else {
				require.Equal(t, user, expUser)
			}
		})
	}
}

func TestUser_DeleteCacheUser(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	type args struct {
		key string
	}
	expUser := TestGetUser(t)
	logging := logger.NewAPILogger("debug")
	testCases := []struct {
		name           string
		wantCreateUser bool
		wantErr        bool
		args           args
		repo           *Repository
		err            error
	}{
		{
			"valid",
			true,
			false,
			args{
				key: expUser.ID.String(),
			},
			NewRepository(rdb, dbx, logging),
			nil,
		},
		{
			"invalid rdb is nil",
			false,
			true,
			args{
				key: expUser.ID.String(),
			},
			NewRepository(nil, nil, logging),
			ErrNotConnection,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.repo.User.DeleteCacheUser(ctx, tc.args.key)
			if err != nil {
				if tc.wantErr {
					require.Equal(t, err, tc.err)
				} else {
					require.Error(t, err)
				}
			}
		})
	}
}

func TestUser_CreateSession(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	type args struct {
		sess   Session
		expire int
	}
	expUser := TestGetUser(t)
	logging := logger.NewAPILogger("debug")
	testCases := []struct {
		name           string
		wantCreateUser bool
		wantErr        bool
		args           args
		repo           *Repository
		err            error
	}{
		{
			"valid",
			true,
			false,
			args{
				sess: Session{
					UserID: expUser.ID,
				},
				expire: 0,
			},
			NewRepository(rdb, dbx, logging),
			nil,
		},
		{
			"invalid rdb is nil",
			false,
			true,
			args{
				sess: Session{
					UserID: expUser.ID,
				},
				expire: 0,
			},
			NewRepository(nil, nil, logging),
			ErrNotConnection,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sessID, err := tc.repo.User.CreateSession(ctx, &tc.args.sess, tc.args.expire)
			if err != nil {
				if tc.wantErr {
					require.Equal(t, err, tc.err)
				} else {
					require.Error(t, err)
				}
			} else {
				require.NotEmpty(t, sessID)
			}
		})
	}
}

func TestUser_GetSessionByID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	type args struct {
		sess   Session
		expire int
	}
	expUser := TestGetUser2(t)
	logging := logger.NewAPILogger("debug")
	testCases := []struct {
		name              string
		wantCreateSession bool
		wantErr           bool
		args              args
		repo              *Repository
		err               error
	}{
		{
			"valid",
			true,
			false,
			args{
				sess: Session{
					UserID: expUser.ID,
				},
				expire: 0,
			},
			NewRepository(rdb, dbx, logging),
			nil,
		},
		{
			"invalid rdb is nil",
			false,
			true,
			args{
				sess: Session{
					UserID: expUser.ID,
				},
				expire: 0,
			},
			NewRepository(nil, nil, logging),
			ErrNotConnection,
		},
		{
			"invalid rdb sessionId nof found",
			false,
			true,
			args{
				sess: Session{
					UserID: expUser.ID,
				},
				expire: 0,
			},
			NewRepository(rdb, dbx, logging),
			grpc_errors.ErrNotFound,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var sessID string = "test"
			var err error
			if tc.wantCreateSession {
				sessID, err = tc.repo.User.CreateSession(ctx, &tc.args.sess, tc.args.expire)
				require.NoError(t, err)
			}

			sess, err := tc.repo.User.GetSessionByID(ctx, sessID)
			if err != nil {
				if tc.wantErr {
					require.Equal(t, err, tc.err)
				} else {
					require.Error(t, err)
				}
			} else {
				require.Equal(t, sess.SessionID, sessID)
				require.Equal(t, sess.UserID, tc.args.sess.UserID)
			}
		})
	}
}

func TestUser_DeleteSessionByID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	type args struct {
		sess   Session
		expire int
	}
	expUser := TestGetUser3(t)
	logging := logger.NewAPILogger("debug")
	testCases := []struct {
		name              string
		wantCreateSession bool
		wantErr           bool
		args              args
		repo              *Repository
		err               error
	}{
		{
			"valid",
			true,
			false,
			args{
				sess: Session{
					UserID: expUser.ID,
				},
				expire: 0,
			},
			NewRepository(rdb, dbx, logging),
			nil,
		},
		{
			"invalid rdb is nil",
			false,
			true,
			args{
				sess: Session{
					UserID: expUser.ID,
				},
				expire: 0,
			},
			NewRepository(nil, nil, logging),
			ErrNotConnection,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var sessID string
			var err error
			if tc.wantCreateSession {
				sessID, err = tc.repo.User.CreateSession(ctx, &tc.args.sess, tc.args.expire)
				require.NoError(t, err)
			}

			err = tc.repo.User.DeleteSessionByID(ctx, sessID)
			if err != nil {
				if tc.wantErr {
					require.Equal(t, err, tc.err)
				} else {
					require.Error(t, err)
				}
			}
		})
	}
}
