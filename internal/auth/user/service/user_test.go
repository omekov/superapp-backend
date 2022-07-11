package service_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/omekov/superapp-backend/internal/auth/domain"
	"github.com/omekov/superapp-backend/internal/auth/user/repository"
	mocksrepository "github.com/omekov/superapp-backend/internal/auth/user/repository/mocks"
	"github.com/omekov/superapp-backend/internal/auth/user/service"
	"github.com/omekov/superapp-backend/pkg/grpc_errors"
	"github.com/omekov/superapp-backend/pkg/jwt"
	mocksjwt "github.com/omekov/superapp-backend/pkg/jwt/mocks"
	"github.com/omekov/superapp-backend/pkg/logger"
	"github.com/omekov/superapp-backend/pkg/mailer"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func mockUserService(t *testing.T) (*service.UserService, *mocksrepository.MockUserer, *mocksjwt.MockJSONWebTokener) {
	t.Helper()

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	userRepo := mocksrepository.NewMockUserer(mockCtl)
	mockJWT := mocksjwt.NewMockJSONWebTokener(mockCtl)
	logg := logger.NewAPILogger("info")
	logg.InitLogger()

	userService := service.NewUserService(userRepo, mockJWT, logg, mailer.Mailer{})
	return userService, userRepo, mockJWT
}

func TestUser_Login(t *testing.T) {
	userService, userRepo, mockJWT := mockUserService(t)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	userID := uuid.New()
	testCases := []struct {
		name         string
		wantErr      bool
		mockBehavior func(r *mocksrepository.MockUserer, jwt *mocksjwt.MockJSONWebTokener, username, password string)
		user         domain.User
		err          error
	}{
		{
			"valid",
			false,
			func(r *mocksrepository.MockUserer, mjwt *mocksjwt.MockJSONWebTokener, username, password string) {
				user := repository.User{
					ID:       userID,
					UserName: "superadmin",
					Password: "$2a$10$TZ5YyQBfHG9t4vG5dhFXreWXx7kGMUy.k.PS11bmAyx6.xySpcwgq",
					State:    service.UserStateEnabled,
				}
				sessioID := "be1f0c2b-c29e-4ba9-a2d3-439fb085992e"
				token := jwt.Token{AccessToken: "token", RefreshToken: "token"}
				r.EXPECT().GetByUsername(ctx, username).Return(user, nil)
				mjwt.EXPECT().PasswordsMatch(user.Password, "password").Return(nil)
				r.EXPECT().CreateSession(ctx, &repository.Session{UserID: user.ID}, 0).Return(sessioID, nil)
				user.Password = ""
				r.EXPECT().SetCacheUser(ctx, sessioID, 0, &user).Return(nil)
				mjwt.EXPECT().NewToken(sessioID).Return(token, nil)
			},
			domain.User{
				Username: "superadmin",
				Password: "password",
			},
			nil,
		},
		{
			"invalid username notfound",
			false,
			func(r *mocksrepository.MockUserer, mjwt *mocksjwt.MockJSONWebTokener, username, password string) {
				r.EXPECT().GetByUsername(ctx, username).Return(repository.User{}, sql.ErrNoRows)

			},
			domain.User{},
			sql.ErrNoRows,
		},
		{
			"invalid user is not activate",
			false,
			func(r *mocksrepository.MockUserer, mjwt *mocksjwt.MockJSONWebTokener, username, password string) {
				user := repository.User{
					UserName: "superadmin",
					Password: "$2a$10$TZ5YyQBfHG9t4vG5dhFXreWXx7kGMUy.k.PS11bmAyx6.xySpcwgq",
					State:    service.UserStateNotActivated,
				}
				r.EXPECT().GetByUsername(ctx, username).Return(user, nil)
			},
			domain.User{
				Username: "superadmin",
				Password: "password",
			},
			service.ErrUserNotActivated,
		},
		{
			"invalid user is not password match",
			false,
			func(r *mocksrepository.MockUserer, mjwt *mocksjwt.MockJSONWebTokener, username, password string) {
				user := repository.User{
					UserName: "superadmin",
					Password: "$2a$10$TZ5YyQBfHG9t4vG5dhFXreWXx7kGMUy.k.PS11bmAyx6.xySpcwgq",
					State:    service.UserStateEnabled,
				}
				r.EXPECT().GetByUsername(ctx, username).Return(user, nil)
				mjwt.EXPECT().PasswordsMatch(user.Password, "password2").Return(bcrypt.ErrMismatchedHashAndPassword)

			},
			domain.User{
				Username: "superadmin",
				Password: "password2",
			},
			bcrypt.ErrMismatchedHashAndPassword,
		},
		{
			"invalid create session failed",
			true,
			func(r *mocksrepository.MockUserer, mjwt *mocksjwt.MockJSONWebTokener, username, password string) {
				user := repository.User{
					ID:       uuid.New(),
					UserName: "superadmin",
					Password: "$2a$10$TZ5YyQBfHG9t4vG5dhFXreWXx7kGMUy.k.PS11bmAyx6.xySpcwgq",
					State:    service.UserStateEnabled,
				}
				r.EXPECT().GetByUsername(ctx, username).Return(user, nil)
				mjwt.EXPECT().PasswordsMatch(user.Password, "password").Return(nil)
				r.EXPECT().CreateSession(ctx, &repository.Session{UserID: user.ID}, 0).Return("", redis.ErrClosed)

			},
			domain.User{
				Username: "superadmin",
				Password: "password",
			},
			redis.ErrClosed,
		},
		{
			"invalid set cache user failed",
			true,
			func(r *mocksrepository.MockUserer, mjwt *mocksjwt.MockJSONWebTokener, username, password string) {
				user := repository.User{
					UserName: "superadmin",
					Password: "$2a$10$TZ5YyQBfHG9t4vG5dhFXreWXx7kGMUy.k.PS11bmAyx6.xySpcwgq",
					State:    service.UserStateEnabled,
				}
				sessionID := "be1f0c2b-c29e-4ba9-a2d3-439fb085992e"

				r.EXPECT().GetByUsername(ctx, username).Return(user, nil)
				mjwt.EXPECT().PasswordsMatch(user.Password, "password").Return(nil)
				r.EXPECT().CreateSession(ctx, &repository.Session{UserID: user.ID}, 0).Return(sessionID, nil)
				user.Password = ""
				r.EXPECT().SetCacheUser(ctx, sessionID, 0, &user).Return(redis.ErrClosed)

			},
			domain.User{
				Username: "superadmin",
				Password: "password",
			},
			redis.ErrClosed,
		},
		{
			"invalid generate token",
			true,
			func(r *mocksrepository.MockUserer, mjwt *mocksjwt.MockJSONWebTokener, username, password string) {
				user := repository.User{
					UserName: "superadmin",
					Password: "$2a$10$TZ5YyQBfHG9t4vG5dhFXreWXx7kGMUy.k.PS11bmAyx6.xySpcwgq",
					State:    service.UserStateEnabled,
				}
				sessionID := "be1f0c2b-c29e-4ba9-a2d3-439fb085992e"

				r.EXPECT().GetByUsername(ctx, username).Return(user, nil)
				mjwt.EXPECT().PasswordsMatch(user.Password, "password").Return(nil)
				r.EXPECT().CreateSession(ctx, &repository.Session{UserID: user.ID}, 0).Return(sessionID, nil)
				user.Password = ""
				r.EXPECT().SetCacheUser(ctx, sessionID, 0, &user).Return(nil)
				mjwt.EXPECT().NewToken(sessionID).Return(jwt.Token{}, jwt.ErrInvalidAccessToken)
			},
			domain.User{
				Username: "superadmin",
				Password: "password",
			},
			jwt.ErrInvalidAccessToken,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior(userRepo, mockJWT, tc.user.Username, tc.user.Password)
			token, err := userService.Login(ctx, tc.user.Username, tc.user.Password)
			if err != nil {
				require.Equal(t, err, tc.err)
			} else {
				tkn := domain.Token{AccessToken: "token", RefreshToken: "token"}
				require.Equal(t, token, tkn)
			}
		})
	}
}

func TestUser_GetMe(t *testing.T) {
	userService, userRepo, mockJWT := mockUserService(t)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	type args struct {
		id uuid.UUID
	}
	userID := uuid.New()
	testCases := []struct {
		name         string
		wantErr      bool
		mockBehavior func(r *mocksrepository.MockUserer, jwt *mocksjwt.MockJSONWebTokener, userID uuid.UUID)
		args         args
		err          error
	}{
		{
			"valid",
			false,
			func(r *mocksrepository.MockUserer, mjwt *mocksjwt.MockJSONWebTokener, userID uuid.UUID) {
				user := repository.User{
					ID:       userID,
					UserName: "superadmin",
					Password: "$2a$10$TZ5YyQBfHG9t4vG5dhFXreWXx7kGMUy.k.PS11bmAyx6.xySpcwgq",
					Email:    "superadmin@mail.kz",
				}
				r.EXPECT().GetCacheByID(ctx, userID.String()).Return(user, nil)
			},
			args{
				id: userID,
			},
			nil,
		},
		{
			"invalid cache user empty",
			false,
			func(r *mocksrepository.MockUserer, mjwt *mocksjwt.MockJSONWebTokener, userID uuid.UUID) {
				user := repository.User{
					ID:       userID,
					UserName: "superadmin",
					Password: "$2a$10$TZ5YyQBfHG9t4vG5dhFXreWXx7kGMUy.k.PS11bmAyx6.xySpcwgq",
					Email:    "superadmin@mail.kz",
				}
				r.EXPECT().GetCacheByID(ctx, userID.String()).Return(repository.User{}, grpc_errors.ErrNotFound)
				r.EXPECT().GetByID(ctx, userID.String()).Return(user, nil)
				r.EXPECT().SetCacheUser(ctx, userID.String(), 0, &user).Return(nil)
			},
			args{
				id: userID,
			},
			nil,
		},
		{
			"invalid user not found",
			false,
			func(r *mocksrepository.MockUserer, mjwt *mocksjwt.MockJSONWebTokener, userID uuid.UUID) {
				user := repository.User{
					ID:       userID,
					UserName: "superadmin",
					Password: "$2a$10$TZ5YyQBfHG9t4vG5dhFXreWXx7kGMUy.k.PS11bmAyx6.xySpcwgq",
					Email:    "superadmin@mail.kz",
				}
				r.EXPECT().GetCacheByID(ctx, userID.String()).Return(repository.User{}, grpc_errors.ErrNotFound)
				r.EXPECT().GetByID(ctx, userID.String()).Return(user, sql.ErrNoRows)
			},
			args{
				id: userID,
			},
			sql.ErrNoRows,
		},
		{
			"invalid set cache user",
			false,
			func(r *mocksrepository.MockUserer, mjwt *mocksjwt.MockJSONWebTokener, userID uuid.UUID) {
				user := repository.User{
					ID:       userID,
					UserName: "superadmin",
					Password: "$2a$10$TZ5YyQBfHG9t4vG5dhFXreWXx7kGMUy.k.PS11bmAyx6.xySpcwgq",
					Email:    "superadmin@mail.kz",
				}
				r.EXPECT().GetCacheByID(ctx, userID.String()).Return(repository.User{}, grpc_errors.ErrNotFound)
				r.EXPECT().GetByID(ctx, userID.String()).Return(user, nil)
				r.EXPECT().SetCacheUser(ctx, userID.String(), 0, &user).Return(redis.ErrClosed)

			},
			args{
				id: userID,
			},
			redis.ErrClosed,
		},
		{
			"invalid get cache by id",
			false,
			func(r *mocksrepository.MockUserer, mjwt *mocksjwt.MockJSONWebTokener, userID uuid.UUID) {
				r.EXPECT().GetCacheByID(ctx, userID.String()).Return(repository.User{}, redis.ErrClosed)
			},
			args{
				id: userID,
			},
			redis.ErrClosed,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior(userRepo, mockJWT, tc.args.id)
			user, err := userService.GetMe(ctx, tc.args.id.String())
			if err != nil {
				require.Equal(t, err, tc.err)
			} else {
				u := domain.User{
					ID:       userID,
					Username: "superadmin",
					Email:    "superadmin@mail.kz",
				}
				require.Equal(t, user, u)
			}
		})
	}
}

func TestRegister(t *testing.T) {
	userService, userRepo, mockJWT := mockUserService(t)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	type args struct {
		user domain.User
	}
	userID := uuid.New()
	testCases := []struct {
		name         string
		wantErr      bool
		mockBehavior func(r *mocksrepository.MockUserer, jwt *mocksjwt.MockJSONWebTokener, user domain.User)
		args         args
		err          error
	}{
		{
			"valid",
			false,
			func(r *mocksrepository.MockUserer, mjwt *mocksjwt.MockJSONWebTokener, user domain.User) {
				passwordHash := "$2a$10$TZ5YyQBfHG9t4vG5dhFXreWXx7kGMUy.k.PS11bmAyx6.xySpcwgq"
				repoUser := repository.User{
					ID:       user.ID,
					UserName: user.Username,
					Password: passwordHash,
					Email:    user.Email,
					State:    service.UserStateNotActivated,
					PinCode:  "123456",
				}
				r.EXPECT().GetByEmail(ctx, user.Email).Return(repository.User{}, nil)
				r.EXPECT().GetByUsername(ctx, user.Username).Return(repository.User{}, nil)
				mjwt.EXPECT().Hash(user.Password).Return(passwordHash, nil)
				r.EXPECT().Create(ctx, repoUser).Return(uuid.Nil, nil)
			},
			args{
				user: domain.User{
					ID:       userID,
					Username: "superadmin2",
					Password: "password",
					Email:    "superadmin2@mail.kz",
				},
			},
			nil,
		},
		{
			"invalid user is already exits",
			false,
			func(r *mocksrepository.MockUserer, mjwt *mocksjwt.MockJSONWebTokener, user domain.User) {
				userData := repository.User{
					ID:       userID,
					UserName: "superadmin",
					Password: "$2a$10$TZ5YyQBfHG9t4vG5dhFXreWXx7kGMUy.k.PS11bmAyx6.xySpcwgq",
					Email:    "superadmin@mail.kz",
				}
				r.EXPECT().GetByEmail(ctx, user.Email).Return(userData, nil)
			},
			args{
				user: domain.User{
					ID:       userID,
					Username: "superadmin",
					Password: "$2a$10$TZ5YyQBfHG9t4vG5dhFXreWXx7kGMUy.k.PS11bmAyx6.xySpcwgq",
					Email:    "superadmin@mail.kz",
				},
			},
			service.ErrUserAlreadyExits,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior(userRepo, mockJWT, tc.args.user)
			err := userService.Register(ctx, tc.args.user)
			if err != nil {
				require.Equal(t, err, tc.err)
			}
		})
	}
}
