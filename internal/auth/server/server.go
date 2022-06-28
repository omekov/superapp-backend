package server

import (
	"context"
	"crypto/tls"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/omekov/superapp-backend/internal/auth/config"
	mygrpc "github.com/omekov/superapp-backend/internal/auth/delivery/grpc"
	"github.com/omekov/superapp-backend/internal/auth/delivery/grpc/v1/proto"
	"github.com/omekov/superapp-backend/internal/auth/interceptors"
	"github.com/omekov/superapp-backend/internal/auth/user/repository"
	"github.com/omekov/superapp-backend/internal/auth/user/service"
	"github.com/omekov/superapp-backend/pkg/conn"
	"github.com/omekov/superapp-backend/pkg/jwt"
	"github.com/omekov/superapp-backend/pkg/logger"
	"github.com/omekov/superapp-backend/pkg/mailer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

// добавить при завершения приложения вытащить logg
// Run ...
func Run(port, cfgPath string) error {

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	logg := logger.NewAPILogger("info")
	logg.InitLogger()

	_cfg := config.New(logg)
	cfg, err := _cfg.Get(cfgPath)
	if err != nil {
		return err
	}

	connect := conn.New(logg)

	dbx := connect.SQLXConn(ctx, cfgPath)
	// if err := goose.Up(dbx.DB, "../../migrations/auth", goose.WithAllowMissing()); err != nil {
	// 	return err
	// }

	rdb := connect.RedisConn(ctx, cfgPath)

	repo := repository.NewRepository(rdb, dbx, logg)

	jwt := jwt.New([]byte(cfg.JWT.Access), []byte(cfg.JWT.Refresh), []byte(cfg.JWT.Mail), 5, 15, 1440)

	mail := mailer.New(cfg.Mailer)

	serv := service.NewService(repo, jwt, logg, mail)
	im := interceptors.NewInterceptorManager(logg)

	opts := []grpc.ServerOption{
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: cfg.GRPC.MaxConnectionIdle * time.Minute,
			Timeout:           cfg.GRPC.Timeout * time.Second,
			MaxConnectionAge:  cfg.GRPC.MaxConnectionAge * time.Minute,
			Time:              cfg.GRPC.Time * time.Minute,
		}),
		grpc.UnaryInterceptor(im.Logger),
		grpc.ChainUnaryInterceptor(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpcrecovery.UnaryServerInterceptor(),
		),
	}

	if cfg.GRPC.TLS {
		cert, err := tls.LoadX509KeyPair(cfg.GRPC.CertFile, cfg.GRPC.KeyFile)
		if err != nil {
			logg.Fatal("tls LoadX509KeyPair", err.Error())
		}
		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{cert},
			ClientAuth:   tls.NoClientCert,
		}
		opts = append(opts, grpc.Creds(credentials.NewTLS(tlsConfig)))
	}

	grpcServer := grpc.NewServer(opts...)
	proto.RegisterAuthServer(grpcServer, &mygrpc.Server{Service: serv, Logg: logg})

	go func() {

		l, err := net.Listen("tcp", port)
		if err != nil {
			logg.Fatal("tcp connection err: ", err.Error())
		}
		defer l.Close()

		if err := grpcServer.Serve(l); err != nil {
			logg.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	<-quit
	grpcServer.GracefulStop()
	logg.Info("Server Exited Properly")

	return nil
}
