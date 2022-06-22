package server

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	_ "github.com/jackc/pgx/v4/stdlib" // ignore
	"github.com/omekov/superapp-backend/internal/auth/interceptors"
	carrepository "github.com/omekov/superapp-backend/internal/salecar/car/repository"
	"github.com/omekov/superapp-backend/internal/salecar/car/service"
	"github.com/omekov/superapp-backend/pkg/conn"
	"github.com/omekov/superapp-backend/pkg/logger"
	"github.com/pressly/goose/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

const addr = ":4041"

// Run - server
func Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logging := logger.NewAPILogger("info")
	logging.InitLogger()
	connect := conn.New(logging)

	dbx := connect.SQLXConn(ctx, "../../configs/config.yaml")
	if err := goose.Up(dbx.DB, "../../migrations/auth", goose.WithAllowMissing()); err != nil {
		return err
	}

	carrepository := carrepository.NewCarRepository(dbx)
	_ = service.NewService(carrepository)
	im := interceptors.NewInterceptorManager(logging)

	grpcServer := grpc.NewServer(grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionIdle: 5 * time.Minute,
		Timeout:           5 * time.Second,
		MaxConnectionAge:  5 * time.Minute,
		Time:              5 * time.Minute,
	}),
		grpc.UnaryInterceptor(im.Logger),
		grpc.ChainUnaryInterceptor(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpcrecovery.UnaryServerInterceptor(),
		),
	)

	go func() {
		l, err := net.Listen("tcp", addr)
		if err != nil {
			logging.Fatal("tcp connection err: ", err.Error())
		}
		defer l.Close()

		if err := grpcServer.Serve(l); err != nil {
			log.Fatal(err)
		}
	}()

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	<-quit
	grpcServer.GracefulStop()
	logging.Info("Server Exited Properly\n")

	return nil
}
