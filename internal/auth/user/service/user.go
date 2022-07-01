package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/omekov/superapp-backend/internal/auth/domain"
	"github.com/omekov/superapp-backend/internal/auth/user/repository"
	"github.com/omekov/superapp-backend/pkg/grpc_errors"
	"github.com/omekov/superapp-backend/pkg/jwt"
	"github.com/omekov/superapp-backend/pkg/logger"
	"github.com/omekov/superapp-backend/pkg/mailer"
	"github.com/omekov/superapp-backend/pkg/util"
)

type UserService struct {
	userRepository repository.Userer
	jwt            *jwt.JWT
	logg           *logger.APILogger
	mailer         mailer.Mailer
}

const (
	UserStateEnabled      = "enabled"
	UserStateNotActivated = "notactivated"
)

var (
	ErrUserNotActivated     = errors.New("is user account not activate")
	ErrUserPinCodeInCorrect = errors.New("is user account pin code incorrect")
	ErrUserAlreadyExits     = errors.New("is user already exits")
	ErrUserNotFound         = errors.New("user not found")
)

// method reset password
// google apple yandex
func NewUserService(
	repository repository.Userer,
	jwt *jwt.JWT,
	logg *logger.APILogger,
	mailer mailer.Mailer,
) *UserService {
	return &UserService{userRepository: repository, jwt: jwt, logg: logg, mailer: mailer}
}

func (s *UserService) Login(ctx context.Context, username, password string) (domain.Token, error) {
	userData, err := s.userRepository.GetByUsername(ctx, username)
	if err != nil {
		return domain.Token{}, err
	}

	if userData.State == UserStateNotActivated {
		return domain.Token{}, ErrUserNotActivated
	}

	err = s.jwt.PasswordsMatch(userData.Password, password)
	if err != nil {
		return domain.Token{}, err
	}

	sessionID, err := s.userRepository.CreateSession(ctx, &repository.Session{
		UserID: userData.ID,
	}, 0)
	if err != nil {
		return domain.Token{}, err
	}

	err = s.userRepository.SetCacheUser(ctx, userData.ID.String(), 0, &userData)
	if err != nil {
		return domain.Token{}, err
	}

	token, err := s.jwt.NewToken(sessionID)
	if err != nil {
		return domain.Token{}, err
	}

	return domain.Token{
		AccessToken:  token.AccessToken,
		RefreshToken: token.Refreshtoken,
	}, nil
}

func (s *UserService) GetMe(ctx context.Context, id uuid.UUID) (domain.User, error) {
	user, err := s.userRepository.GetCacheByID(ctx, id.String())
	if errors.Is(err, grpc_errors.ErrNotFound) {
		user, err = s.userRepository.GetByID(ctx, id)
		if err != nil {
			return domain.User{}, fmt.Errorf("GetByID: %v", err)
		}

		err = s.userRepository.SetCacheUser(ctx, user.ID.String(), 0, &user)
		if err != nil {
			return domain.User{}, fmt.Errorf("SetCacheUser: %v", err)
		}
		err = nil
	}

	if err != nil {
		return domain.User{}, fmt.Errorf("GetCacheByID: %v", err)
	}

	return domain.User{
		ID:       user.ID,
		Username: user.UserName,
		Password: user.Password,
	}, nil
}

// Отправка на почту код активацию
func (s *UserService) Register(ctx context.Context, user domain.User) error {

	userData, err := s.userRepository.GetByEmail(ctx, user.Email)
	if err != nil {
		return err
	}
	if userData.Email == user.Email {
		return ErrUserAlreadyExits
	}

	userData, err = s.userRepository.GetByUsername(ctx, user.Username)
	if err != nil {
		return err
	}

	if userData.UserName == user.Username {
		return ErrUserAlreadyExits
	}

	password, err := s.jwt.Hash(user.Password)
	if err != nil {
		return err
	}

	pinCode := util.GenPinCode(6)

	repoUser := repository.User{
		ID:       uuid.New(),
		UserName: user.Username,
		Password: password,
		Email:    user.Email,
		State:    UserStateNotActivated,
		PinCode:  pinCode,
	}
	_, err = s.userRepository.Create(ctx, repoUser)
	if err != nil {
		return err
	}

	go func() {
		if err := s.mailer.Send(user.Email, "activate.html", repoUser); err != nil {
			s.logg.Error(err)
			return
		}
	}()

	return nil
}

func (s *UserService) Refresh(ctx context.Context, refToken string) (domain.Token, error) {
	claims, err := s.jwt.GetClaimsRefresh(refToken)
	if err != nil {
		return domain.Token{}, err
	}
	err = claims.Valid()
	if err != nil {
		return domain.Token{}, err
	}

	token, err := s.jwt.NewToken(claims.SessionID)
	if err != nil {
		return domain.Token{}, err
	}

	return domain.Token{
		AccessToken:  token.AccessToken,
		RefreshToken: token.Refreshtoken,
	}, nil
}

func (s *UserService) ForgetPassword(ctx context.Context, email string) error {
	user, err := s.userRepository.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrUserNotFound
		}
		return err
	}

	token, err := s.jwt.NewMailToken(email)
	if err != nil {
		return err
	}

	go func() {
		pass := struct {
			PassToken string
			Username  string
		}{
			PassToken: token,
			Username:  user.UserName,
		}

		if err := s.mailer.Send(email, "forget_password.html", pass); err != nil {
			s.logg.Error(err)
			return
		}
	}()
	return nil
}
func (s *UserService) ResetPassword(ctx context.Context, passToken, newPassword string) error {
	claims, err := s.jwt.GetClaimsMail(passToken)
	if err != nil {
		return err
	}
	email := claims.SessionID

	user, err := s.userRepository.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrUserNotFound
		}
		return err
	}

	password, err := s.jwt.Hash(newPassword)
	if err != nil {
		return err
	}

	err = s.userRepository.UpdatePassword(ctx, user.ID, password)
	if err != nil {
		return err
	}

	go func() {
		if err := s.mailer.Send(email, "success_password.html", user); err != nil {
			s.logg.Error(err)
			return
		}
	}()
	return nil
}

func (s *UserService) Activate(ctx context.Context, email, pinCode string) error {

	userData, err := s.userRepository.GetByEmail(ctx, email)
	if err != nil {
		return err
	}

	if userData.PinCode != pinCode {
		return ErrUserPinCodeInCorrect
	}

	sanitizePinCode := "-"
	err = s.userRepository.UpdateState(ctx, userData.ID, UserStateEnabled, sanitizePinCode)
	if err != nil {
		return err
	}

	return nil
}
