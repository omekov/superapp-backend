package interceptors

import (
	"context"
	"time"

	grpcauth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/omekov/superapp-backend/internal/auth/user/service"
	"github.com/omekov/superapp-backend/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type (
	authenticated      bool
	sessionIDFromClaim string
)

const (
	Authenticated      authenticated      = false
	SessionIDFromClaim sessionIDFromClaim = "sessionID"
)

// InterceptorManager ...
type InterceptorManager struct {
	logger            logger.Logger
	noVerifyMethodMap map[string]bool
	service           service.Service
}

// InterceptorManager ...
func NewInterceptorManager(logger logger.Logger, service service.Service, noVerifyMethodMap map[string]bool) *InterceptorManager {
	return &InterceptorManager{logger: logger, service: service, noVerifyMethodMap: noVerifyMethodMap}
}

// Logger ...
func (im *InterceptorManager) Logger(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	start := time.Now()
	md, _ := metadata.FromIncomingContext(ctx)
	reply, err := handler(ctx, req)
	im.logger.Infof("Method: %s, Time: %v, Metadata: %v, Err: %v", info.FullMethod, time.Since(start), md, err)

	return reply, err
}

func (im *InterceptorManager) AuthFunc(ctx context.Context) (context.Context, error) {
	methodString, has := grpc.Method(ctx)
	if has {
		if value, ok := im.noVerifyMethodMap[methodString]; ok && value {
			return ctx, nil
		}
	}

	token, err := grpcauth.AuthFromMD(ctx, "Bearer")
	if err != nil {
		return nil, err
	}

	sessionID, err := im.service.User.VerifyToken(ctx, token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid auth token: %v", err)
	}

	// set meta data in context for possible later validation
	newCtx := context.WithValue(ctx, Authenticated, true)
	newCtx = context.WithValue(newCtx, SessionIDFromClaim, sessionID)

	return newCtx, nil
}
