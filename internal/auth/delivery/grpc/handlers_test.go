package grpc_test

/*
import (
	"context"
	"log"
	"net"
	"testing"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	. "github.com/omekov/superapp-backend/internal/auth/delivery/grpc"
	"github.com/omekov/superapp-backend/internal/auth/delivery/grpc/v1/proto"
	"github.com/stretchr/testify/require"
	gogrpc "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

func dialer() func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1024 * 1024)
	server := gogrpc.NewServer()
	proto.RegisterAuthServer(server, &Server{})

	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()
	return func(ctx context.Context, s string) (net.Conn, error) {
		return listener.Dial()
	}
}

func TestAuthServer_Login(t *testing.T) {
	testCases := []struct {
		name    string
		res     *proto.AuthResponse
		errCode codes.Code
		errMsg  string
	}{
		{
			"valid",
			&proto.AuthResponse{AccessToken: "", RefreshToken: ""},
			codes.OK,
			validation.ErrLengthInvalid.Error(),
		},
	}

	ctx := context.Background()
	conn, err := gogrpc.DialContext(ctx, "", gogrpc.WithInsecure(), gogrpc.WithContextDialer(dialer()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := proto.NewAuthClient(conn)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := &proto.AuthRequest{
				Username: "azamat",
				Password: "azamat12",
			}

			res, err := client.Login(ctx, req)
			if res != nil {
				require.Equal(t, tc.res.AccessToken, res.AccessToken)
				require.Equal(t, tc.res.RefreshToken, res.RefreshToken)
			}

			if err != nil {
				if err, ok := status.FromError(err); ok {
					// if err.Code() != tc.errCode {
					// 	t.Error("error code: expected", codes.InvalidArgument, "received", err.Code())
					// }
					if err.Message() != tc.errMsg {
						t.Error("error message: expected", tc.errMsg, "received", err.Message())
					}
				}
			}

		})
	}
}
*/
