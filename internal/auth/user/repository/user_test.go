package repository_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/google/uuid"
	. "github.com/omekov/superapp-backend/internal/auth/user/repository"
	"github.com/omekov/superapp-backend/pkg/logger"
	"github.com/stretchr/testify/require"
)

func TestUser_Create(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	logging := logger.NewAPILogger("debug")
	testCases := []struct {
		name       string
		wantErr    bool
		createUser User
		repo       *Repository
	}{
		{
			"valid",
			false,
			TestGetUser(t),
			NewRepository(rdb, dbx, logging),
		},
		{
			"invalid dbx is nil",
			true,
			TestGetUser(t),
			NewRepository(rdb, nil, logging),
		},
		{
			"invalid duplicated err",
			false,
			TestGetUser(t),
			NewRepository(rdb, dbx, logging),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			userID, err := tc.repo.User.Create(ctx, tc.createUser)
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

func TestUser_GetByUsername(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	logging := logger.NewAPILogger("debug")
	testCases := []struct {
		name           string
		wantCreateUser bool
		wantErr        bool
		payload        User
		repo           *Repository
		err            error
	}{
		{
			"valid",
			true,
			false,
			TestGetUser2(t),
			NewRepository(rdb, dbx, logging),
			nil,
		},
		{
			"invalid dbx is nil",
			false,
			true,
			TestGetUser2(t),
			NewRepository(rdb, nil, logging),
			ErrNotConnection,
		},
		{
			"invalid no rows",
			false,
			true,
			User{},
			NewRepository(rdb, dbx, logging),
			sql.ErrNoRows,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.wantCreateUser {
				_, err := tc.repo.User.Create(ctx, tc.payload)
				require.NoError(t, err)
			}

			user, err := tc.repo.User.GetByUsername(ctx, tc.payload.UserName)
			if err != nil {
				if tc.wantErr {
					require.Equal(t, err, tc.err)
				} else {
					require.Error(t, err)
				}
				return
			}
			require.Equal(t, tc.payload, user)
		})
	}
}

func TestUser_GetByID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	logging := logger.NewAPILogger("debug")
	testCases := []struct {
		name           string
		wantCreateUser bool
		wantErr        bool
		payload        User
		repo           *Repository
		err            error
	}{
		{
			"valid",
			true,
			false,
			TestGetUser3(t),
			NewRepository(rdb, dbx, logging),
			nil,
		},
		{
			"invalid dbx is nil",
			false,
			true,
			TestGetUser3(t),
			NewRepository(rdb, nil, logging),
			ErrNotConnection,
		},
		{
			"invalid no rows",
			false,
			true,
			User{},
			NewRepository(rdb, dbx, logging),
			sql.ErrNoRows,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.wantCreateUser {
				_, err := tc.repo.User.Create(ctx, tc.payload)
				require.NoError(t, err)
			}

			user, err := tc.repo.User.GetByID(ctx, tc.payload.ID.String())
			if err != nil {
				if tc.wantErr {
					require.Equal(t, err, tc.err)
				} else {
					require.Error(t, err)
				}
				return
			}
			require.Equal(t, tc.payload, user)
		})
	}
}

func TestUser_GetByEmail(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	logging := logger.NewAPILogger("debug")
	testCases := []struct {
		name           string
		wantCreateUser bool
		wantErr        bool
		payload        User
		repo           *Repository
		err            error
	}{
		{
			"valid",
			true,
			false,
			TestGetUser4(t),
			NewRepository(rdb, dbx, logging),
			nil,
		},
		{
			"invalid dbx is nil",
			false,
			true,
			TestGetUser4(t),
			NewRepository(rdb, nil, logging),
			ErrNotConnection,
		},
		{
			"invalid no rows",
			false,
			true,
			User{},
			NewRepository(rdb, dbx, logging),
			sql.ErrNoRows,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.wantCreateUser {
				_, err := tc.repo.User.Create(ctx, tc.payload)
				require.NoError(t, err)
			}

			user, err := tc.repo.User.GetByEmail(ctx, tc.payload.Email)
			if err != nil {
				if tc.wantErr {
					require.Equal(t, err, tc.err)
				} else {
					require.Error(t, err)
				}
				return
			}
			require.Equal(t, tc.payload, user)
		})
	}
}

func TestUser_UpdatePassword(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	logging := logger.NewAPILogger("debug")
	newPassword := "test"
	testCases := []struct {
		name           string
		wantCreateUser bool
		wantErr        bool
		payload        User
		repo           *Repository
		err            error
	}{
		{
			"valid",
			true,
			false,
			TestGetUser5(t),
			NewRepository(rdb, dbx, logging),
			nil,
		},
		{
			"invalid dbx is nil",
			false,
			true,
			TestGetUser5(t),
			NewRepository(rdb, nil, logging),
			ErrNotConnection,
		},
		{
			"invalid no rows",
			false,
			true,
			User{
				ID: uuid.New(),
			},
			NewRepository(rdb, dbx, logging),
			ErrNoRowsUpdated,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.wantCreateUser {
				_, err := tc.repo.User.Create(ctx, tc.payload)
				require.NoError(t, err)
			}

			err := tc.repo.User.UpdatePassword(ctx, tc.payload.ID, newPassword)
			if err != nil {
				if tc.wantErr {
					require.Equal(t, err, tc.err)
				} else {
					require.Error(t, err)
				}
			}

			user, err := tc.repo.User.GetByID(ctx, tc.payload.ID.String())
			if err != nil {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, user.Password, newPassword)
			}
		})
	}
}

func Test_UpdateState(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	type args struct {
		pinCode string
		state   string
	}
	logging := logger.NewAPILogger("debug")
	testCases := []struct {
		name           string
		wantCreateUser bool
		wantErr        bool
		payload        User
		args           args
		repo           *Repository
		err            error
	}{
		{
			"valid",
			true,
			false,
			TestGetUser6(t),
			args{
				pinCode: "123456",
				state:   "enabled",
			},
			NewRepository(rdb, dbx, logging),
			nil,
		},
		{
			"invalid dbx is nil",
			false,
			true,
			TestGetUser6(t),
			args{
				pinCode: "123456",
				state:   "enabled",
			},
			NewRepository(rdb, nil, logging),
			ErrNotConnection,
		},
		{
			"invalid no rows",
			false,
			true,
			User{
				ID: uuid.New(),
			},
			args{
				pinCode: "123456",
				state:   "enabled",
			},
			NewRepository(rdb, dbx, logging),
			ErrNoRowsUpdated,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.wantCreateUser {
				_, err := tc.repo.User.Create(ctx, tc.payload)
				require.NoError(t, err)
			}

			err := tc.repo.User.UpdateState(ctx, tc.payload.ID, tc.args.state, tc.args.pinCode)
			if err != nil {
				if tc.wantErr {
					require.Equal(t, err, tc.err)
				} else {
					require.Error(t, err)
				}
			}

			user, err := tc.repo.User.GetByEmail(ctx, tc.payload.Email)
			if err != nil {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, user.PinCode, tc.args.pinCode)
				require.Equal(t, user.State, tc.args.state)
			}
		})
	}
}
