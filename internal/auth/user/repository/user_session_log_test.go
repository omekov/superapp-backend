package repository_test

import (
	"context"
	"testing"
	"time"

	. "github.com/omekov/superapp-backend/internal/auth/user/repository"
	"github.com/omekov/superapp-backend/pkg/logger"
	"github.com/stretchr/testify/require"
)

func TestCreateUserSessionLog(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	logging := logger.NewAPILogger("debug")
	testCases := []struct {
		name                 string
		wantErr              bool
		createUserSessionLog UserSessionLog
		repo                 *Repository
	}{
		{
			"valid",
			false,
			UserSessionLog{
				SessionID:   "30d45853-0ca2-44fa-9e3c-a708f67a1dbb",
				Username:    "test",
				UserAgent:   "MacOS",
				ClientIP:    "127.0.0.1",
				HTTPMethod:  "POST",
				HTTPPath:    "/login",
				HTTPReqBody: "",
			},
			NewRepository(rdb, dbx, logging),
		},
		{
			"invalid dbx is nil",
			true,
			UserSessionLog{},
			NewRepository(rdb, nil, logging),
		},
		{
			"invalid duplicated err",
			false,
			UserSessionLog{},
			NewRepository(rdb, dbx, logging),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			userID, err := tc.repo.User.CreateUserSessionLog(ctx, tc.createUserSessionLog)
			if err != nil {
				if tc.wantErr {
					require.Equal(t, err, ErrNotConnection)
				} else {
					require.Error(t, err)
				}
				return
			}
			require.NoError(t, err)
			require.NotEmpty(t, userID)
		})
	}
}

func TestUpdateUserSessionLog(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	logging := logger.NewAPILogger("debug")
	testCases := []struct {
		name             string
		wantErr          bool
		userSessionLogID uint
		userSessionLog   UserSessionLog
		repo             *Repository
	}{
		{
			"valid",
			false,
			1,
			UserSessionLog{
				SessionID:   "30d45853-0ca2-44fa-9e3c-a708f67a1dbb",
				Username:    "test",
				UserAgent:   "MacOS",
				ClientIP:    "127.0.0.1",
				HTTPMethod:  "POST",
				HTTPPath:    "/login",
				HTTPReqBody: "",
			},
			NewRepository(rdb, dbx, logging),
		},
		{
			"invalid dbx is nil",
			true,
			1,
			UserSessionLog{},
			NewRepository(rdb, nil, logging),
		},
		{
			"invalid duplicated err",
			false,
			1,
			UserSessionLog{},
			NewRepository(rdb, dbx, logging),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			userSessionLogID, _ := tc.repo.User.CreateUserSessionLog(ctx, tc.userSessionLog)
			err := tc.repo.User.UpdateUserSessionLog(ctx, userSessionLogID, tc.userSessionLog)
			if err != nil {
				if tc.wantErr {
					require.Equal(t, err, ErrNotConnection)
				} else {
					require.Error(t, err)
				}
				return
			}
			require.NoError(t, err)
		})
	}
}
