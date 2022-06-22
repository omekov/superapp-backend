package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/omekov/superapp-backend/internal/auth/domain"
	"github.com/omekov/superapp-backend/internal/auth/user/repository"
	mock_repository "github.com/omekov/superapp-backend/internal/auth/user/repository/mocks"
	"github.com/omekov/superapp-backend/internal/auth/user/service"
	"github.com/omekov/superapp-backend/pkg/jwt"
	"github.com/omekov/superapp-backend/pkg/logger"
	"github.com/omekov/superapp-backend/pkg/mailer"
	"github.com/stretchr/testify/require"
)

var errInternalServErr = errors.New("test: internal server error")

func mockUserService(t *testing.T) (*service.UserService, *mock_repository.MockUserer) {
	t.Helper()

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	userRepo := mock_repository.NewMockUserer(mockCtl)
	logg := logger.NewAPILogger("info")
	logg.InitLogger()
	jwt := jwt.New([]byte("access"), []byte("refresh"), []byte("mail"), 5, 15, 1440)

	userService := service.NewUserService(userRepo, jwt, logg, mailer.Mailer{})
	return userService, userRepo
}

func TestUser_LoginErr(t *testing.T) {
	ctx := context.Background()
	userService, userRepo := mockUserService(t)

	userRepo.EXPECT().GetByUsername(ctx, gomock.Any()).Return(repository.User{}, errInternalServErr)
	// userRepo.EXPECT().CreateSession(ctx, gomock.Any(), gomock.Any()).Return("", errInternalServErr)
	// userRepo.EXPECT().SetCacheUser(ctx, gomock.Any(), gomock.Any(), gomock.Any()).Return(errInternalServErr)

	token, err := userService.Login(ctx, "test", "test")
	require.True(t, errors.Is(err, errInternalServErr))
	require.Equal(t, domain.Token{}, token)
}

func TestUser_Login(t *testing.T) {
	ctx := context.Background()
	userService, userRepo := mockUserService(t)

	userRepo.EXPECT().GetByUsername(ctx, gomock.Any())
	userRepo.EXPECT().CreateSession(ctx, gomock.Any(), gomock.Any())
	userRepo.EXPECT().SetCacheUser(ctx, gomock.Any(), gomock.Any(), gomock.Any())

	token, err := userService.Login(ctx, "test", "test12")
	require.NoError(t, err)
	require.IsType(t, domain.Token{}, token)

}
