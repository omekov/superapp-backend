package interceptors

import (
	"context"
	"time"

	"github.com/omekov/superapp-backend/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// InterceptorManager ...
type InterceptorManager struct {
	logger logger.Logger
}

// InterceptorManager ...
func NewInterceptorManager(logger logger.Logger) *InterceptorManager {
	return &InterceptorManager{logger: logger}
}

// Logger ...
func (im *InterceptorManager) Logger(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	start := time.Now()
	md, _ := metadata.FromIncomingContext(ctx)
	reply, err := handler(ctx, req)
	im.logger.Infof("Method: %s, Time: %v, Metadata: %v, Err: %v", info.FullMethod, time.Since(start), md, err)

	return reply, err
}
