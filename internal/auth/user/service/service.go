package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/omekov/superapp-backend/internal/auth/domain"
	"github.com/omekov/superapp-backend/internal/auth/user/repository"
	"github.com/omekov/superapp-backend/pkg/jwt"
	"github.com/omekov/superapp-backend/pkg/logger"
	"github.com/omekov/superapp-backend/pkg/mailer"
)

type Service struct {
	User UserServicier
}

func NewService(repository *repository.Repository, jwt *jwt.JWT, logg *logger.APILogger, mail mailer.Mailer) Service {
	return Service{
		User: NewUserService(repository.User, jwt, logg, mail),
	}
}

type UserServicier interface {
	Register(ctx context.Context, user domain.User) error
	Login(ctx context.Context, username, password string) (domain.Token, error)
	GetMe(ctx context.Context, id uuid.UUID) (domain.User, error)
	Refresh(ctx context.Context, refToken string) (domain.Token, error)
	ResetPassword(ctx context.Context, passToken, newPassword string) error
	ForgetPassword(ctx context.Context, email string) error
	Activate(ctx context.Context, email, pinCode string) error
}
