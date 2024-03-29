package grpc_errors // ignore

import (
	"context"
	"database/sql"
	"net/http"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
)

var (
	// ErrNotFound ...
	ErrNotFound = errors.New("Not found")
	// ErrNoCtxMetaData ...
	ErrNoCtxMetaData = errors.New("No ctx metadata")
	// ErrInvalidSessionID ...
	ErrInvalidSessionID = errors.New("Invalid session id")
	// ErrEmailExists ...
	ErrEmailExists = errors.New("Email already exists")
)

// ParseGRPCErrStatusCode - error and get code
func ParseGRPCErrStatusCode(err error) codes.Code {
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return codes.NotFound
	case errors.Is(err, redis.Nil):
		return codes.NotFound
	case errors.Is(err, context.Canceled):
		return codes.Canceled
	case errors.Is(err, context.DeadlineExceeded):
		return codes.DeadlineExceeded
	case errors.Is(err, ErrEmailExists):
		return codes.AlreadyExists
	case errors.Is(err, ErrNoCtxMetaData):
		return codes.Unauthenticated
	case errors.Is(err, ErrInvalidSessionID):
		return codes.PermissionDenied
	case strings.Contains(err.Error(), "Validate"):
		return codes.InvalidArgument
	case strings.Contains(err.Error(), "redis"):
		return codes.NotFound
	}
	return codes.Internal
}

// MapGRPCErrCodeToHTTPStatus - errors codes to http status
func MapGRPCErrCodeToHTTPStatus(code codes.Code) int {
	switch code {
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	case codes.AlreadyExists:
		return http.StatusBadRequest
	case codes.NotFound:
		return http.StatusNotFound
	case codes.Internal:
		return http.StatusInternalServerError
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.Canceled:
		return http.StatusRequestTimeout
	case codes.DeadlineExceeded:
		return http.StatusGatewayTimeout
	case codes.InvalidArgument:
		return http.StatusBadRequest
	}
	return http.StatusInternalServerError
}
