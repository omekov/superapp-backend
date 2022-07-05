package grpc

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/go-ozzo/ozzo-validation/is"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	proto "github.com/omekov/superapp-backend/internal/auth/delivery/grpc/v1"
	"github.com/omekov/superapp-backend/internal/auth/domain"
	"github.com/omekov/superapp-backend/internal/auth/user/service"
	"github.com/omekov/superapp-backend/pkg/logger"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	// ErrIncorrectUsernamePassword ...
	ErrIncorrectUsernamePassword = errors.New("incorrect Username or Password")
)

// Server ...
type Server struct {
	Logg    *logger.APILogger
	Service service.Service
	proto.UnimplementedAuthServer
}

// Login ...
func (s *Server) Login(ctx context.Context, in *proto.AuthRequest) (*proto.AuthResponse, error) {
	username := strings.ToLower(in.GetUsername())
	err := validation.Validate(username,
		validation.Required,
		validation.Length(6, 100),
		is.LowerCase,
		is.Alphanumeric,
	)
	if err != nil {
		return nil, fmt.Errorf("username invalid: %w", err)
	}

	err = validation.Validate(in.GetPassword(),
		validation.Required,
		validation.Length(6, 100),
		is.Alphanumeric,
	)
	if err != nil {
		return nil, fmt.Errorf("password invalid: %w", err)
	}

	token, err := s.Service.User.Login(ctx, username, in.GetPassword())
	if err != nil {
		s.Logg.Warnf("Login: %s - %s", username, err.Error())
		return nil, ErrIncorrectUsernamePassword
	}
	return &proto.AuthResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}, err
}

// GetMe ...
func (s *Server) GetMe(ctx context.Context, in *proto.GetMeRequest) (*proto.GetMeResponse, error) {
	user, err := s.Service.User.GetMe(ctx, in.GetSessionID())
	return &proto.GetMeResponse{
		User: &proto.User{
			ID:       user.ID.String(),
			Username: user.Username,
			Email:    user.Email,
		},
	}, err
}

// Register ...
func (s *Server) Register(ctx context.Context, in *proto.UserRequest) (*emptypb.Empty, error) {
	username := strings.ToLower(in.GetUsername())
	err := validation.Validate(username,
		validation.Required,
		validation.Length(6, 100),
		is.LowerCase,
		is.Alphanumeric,
	)
	if err != nil {
		return nil, fmt.Errorf("username invalid: %w", err)
	}

	err = validation.Validate(in.GetPassword(),
		validation.Required,
		validation.Length(8, 100),
		is.Alphanumeric,
	)
	if err != nil {
		return nil, fmt.Errorf("password invalid: %w", err)
	}

	email := strings.ToLower(in.GetEmail())
	err = validation.Validate(email,
		validation.Required,
		is.LowerCase,
	)
	if err != nil {
		return nil, fmt.Errorf("email invalid: %w", err)
	}

	err = s.Service.User.Register(ctx, domain.User{
		Username: username,
		Password: in.GetPassword(),
		Email:    email,
	})

	return &emptypb.Empty{}, err
}

// Refresh ...
func (s *Server) Refresh(ctx context.Context, in *proto.RefreshRequest) (*proto.AuthResponse, error) {
	token, err := s.Service.User.Refresh(ctx, in.GetRefreshToken())
	return &proto.AuthResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}, err
}

// Activate ...
func (s *Server) Activate(ctx context.Context, in *proto.ActivateRequest) (*emptypb.Empty, error) {
	email := strings.ToLower(in.GetEmail())
	err := validation.Validate(email,
		validation.Required,
		is.LowerCase,
	)
	if err != nil {
		return nil, fmt.Errorf("email invalid: %w", err)
	}

	pinCode := strings.ToLower(in.GetPinCode())
	err = validation.Validate(pinCode,
		validation.Required,
		is.LowerCase,
	)
	if err != nil {
		return nil, fmt.Errorf("pinCode invalid: %w", err)
	}

	if err := s.Service.User.Activate(ctx, email, pinCode); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

// ResetPassword ...
func (s *Server) ResetPassword(ctx context.Context, in *proto.ResetPasswordRequest) (*emptypb.Empty, error) {
	err := validation.Validate(in.GetPassToken(),
		validation.Required,
	)
	if err != nil {
		return nil, fmt.Errorf("passToken invalid: %w", err)
	}

	err = validation.Validate(in.GetNewPassword(),
		validation.Required,
	)
	if err != nil {
		return nil, fmt.Errorf("newPassword invalid: %w", err)
	}

	if err := s.Service.User.ResetPassword(ctx, in.GetPassToken(), in.GetNewPassword()); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// ForgetPassword ...
func (s *Server) ForgetPassword(ctx context.Context, in *proto.ForgetPasswordRequest) (*emptypb.Empty, error) {
	email := strings.ToLower(in.GetEmail())
	err := validation.Validate(email,
		validation.Required,
		is.LowerCase,
	)
	if err != nil {
		return nil, fmt.Errorf("email invalid: %w", err)
	}

	if err := s.Service.User.ForgetPassword(ctx, email); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
