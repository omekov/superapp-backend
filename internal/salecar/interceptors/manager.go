package interceptors

import (
	"context"
	"log"
	"time"

	"github.com/omekov/superapp-backend/internal/salecar/car/service"
	"github.com/omekov/superapp-backend/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// InterceptorManager ...
type InterceptorManager struct {
	logger  logger.Logger
	service *service.Service
}

// InterceptorManager ...
func NewInterceptorManager(logger logger.Logger, service *service.Service) *InterceptorManager {
	return &InterceptorManager{logger: logger, service: service}
}

// Logger ...
func (im *InterceptorManager) Logger(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	start := time.Now()
	md, _ := metadata.FromIncomingContext(ctx)
	reply, err := handler(ctx, req)
	im.logger.Infof("Method: %s, Time: %v, Metadata: %v, Err: %v", info.FullMethod, time.Since(start), md, err)

	return reply, err
}

func (im *InterceptorManager) Stream() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		err := im.authorize(stream.Context(), info.FullMethod)
		if err != nil {
			return err
		}
		return handler(srv, stream)
	}
}
func (im *InterceptorManager) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		log.Println("--> unary interceptor: ", info.FullMethod)

		// TODO: implement authorization
		err := im.authorize(ctx, info.FullMethod)
		if err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}

func (im *InterceptorManager) authorize(ctx context.Context, method string) error {
	// accessibleRoles, ok := im.accessibleRoles[method]
	// if !ok {
	// 	// everyone can access
	// 	return nil
	// }

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	values := md["authorization"]
	if len(values) == 0 {
		return status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}

	accessToken := values[0]
	log.Println(accessToken)
	// claims, err := im.jwtManager.Verify(accessToken)
	// if err != nil {
	// 	return status.Errorf(codes.Unauthenticated, "access token is invalid: %v", err)
	// }

	// for _, role := range accessibleRoles {
	// 	if role == claims.Role {
	// 		return nil
	// 	}
	// }

	return status.Error(codes.PermissionDenied, "no permission to access this RPC")
}
