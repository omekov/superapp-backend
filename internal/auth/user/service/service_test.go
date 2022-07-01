package service_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/omekov/superapp-backend/internal/auth/config"
	"github.com/omekov/superapp-backend/internal/auth/user/repository"
	mocksrepository "github.com/omekov/superapp-backend/internal/auth/user/repository/mocks"
	"github.com/omekov/superapp-backend/internal/auth/user/service"
	"github.com/omekov/superapp-backend/pkg/jwt"
	"github.com/omekov/superapp-backend/pkg/logger"
	"github.com/omekov/superapp-backend/pkg/mailer"
	"github.com/stretchr/testify/require"
)

func mockRepository(t *testing.T) *repository.Repository {
	t.Helper()
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	userRepo := mocksrepository.NewMockUserer(mockCtl)
	return &repository.Repository{
		User: userRepo,
	}
}

func TestNewService(t *testing.T) {

	repo := mockRepository(t)
	jwt := jwt.New([]byte("access"), []byte("refresh"), []byte("mail"), 5, 15, 1440)

	logg := logger.NewAPILogger("debug")
	logg.InitLogger()
	mail := mailer.New(config.MailerConfig{})

	want := service.Service{
		User: service.NewUserService(repo.User, jwt, logg, mail),
	}
	got := service.NewService(repo, jwt, logg, mail)
	require.Equal(t, want, got)
}
