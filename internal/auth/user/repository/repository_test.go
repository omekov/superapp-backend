package repository_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/docker/go-connections/nat"
	"github.com/go-redis/redis/v8"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	. "github.com/omekov/superapp-backend/internal/auth/user/repository"
	"github.com/omekov/superapp-backend/pkg/conn"
	"github.com/omekov/superapp-backend/pkg/logger"
	"github.com/pressly/goose/v3"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	dbx *sqlx.DB
	rdb *redis.Client
)

const (
	configPath = "../../../../configs/config.yaml"
)

func TestMain(m *testing.M) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	postgresPort := nat.Port("5432/tcp")
	postgresUser := "postgres"
	postgresPass := "postgres"
	postgresName := "postgres"

	pgContainer, err := testcontainers.GenericContainer(
		ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: testcontainers.ContainerRequest{
				Image:        "postgres:13-alpine",
				ExposedPorts: []string{postgresPort.Port()},
				Env: map[string]string{
					"POSTGRES_DB":       postgresName,
					"POSTGRES_USER":     postgresUser,
					"POSTGRES_PASSWORD": postgresPass,
				},
				WaitingFor: wait.ForAll(
					wait.ForLog("database system is ready to accept connections"),
					wait.ForListeningPort(postgresPort),
				),
			},
			Started: true,
		},
	)
	if err != nil {
		log.Fatal(fmt.Errorf("testcontainers.GenericContainer %s", err))
	}

	defer func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			log.Fatal(fmt.Errorf("pgContrainer.Terminate %s", err))
		}
	}()

	pgContainerPort, err := pgContainer.MappedPort(ctx, "5432")
	if err != nil {
		log.Fatal(fmt.Errorf("pgContrainer.MappedPort %s", err))
	}

	pgContainerHostIP, err := pgContainer.Host(ctx)
	if err != nil {
		log.Fatal(fmt.Errorf("pgContainer.Host %s", err))
	}

	log.Printf("TestContainer Postgres HOST:%s PORT:%s\n", pgContainerHostIP, pgContainerPort.Port())
	logg := logger.NewAPILogger("debug")
	logg.InitLogger()
	connect := conn.New(logg)

	dsn := fmt.Sprintf(
		"host=%s port=%s dbname=%s sslmode=%s user=%s password=%s",
		pgContainerHostIP,
		pgContainerPort.Port(),
		postgresName,
		"disable",
		postgresUser,
		postgresPass,
	)
	dbx, err = sqlx.ConnectContext(ctx, "pgx", dsn)
	if err != nil {
		log.Fatal(fmt.Errorf("sqlx.ConnectContext: %s", err))
	}
	defer dbx.Close()

	err = goose.Up(dbx.DB, "../../../../migrations/auth")
	if err != nil {
		log.Fatal(fmt.Errorf("goose.Up %s", err))
	}

	redisPass := ""
	redisPort := nat.Port("6379/tcp")

	redisContainer, err := testcontainers.GenericContainer(
		ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: testcontainers.ContainerRequest{
				Image:        "redis:6",
				ExposedPorts: []string{redisPort.Port()},
				Env: map[string]string{
					"REDIS_PASSWORD": redisPass,
				},
				WaitingFor: wait.ForAll(
					wait.ForLog("* Ready to accept connections"),
					wait.ForListeningPort(redisPort),
				),
			},
			Started: true,
		})
	if err != nil {
		log.Fatal(fmt.Errorf("testcontainers.GenericContainer %s", err))
	}

	redisContainerPort, err := redisContainer.MappedPort(ctx, "6379")
	if err != nil {
		log.Fatal(fmt.Errorf("redisContainer.MappedPort %s", err))
	}

	hostIP, err := redisContainer.Host(ctx)
	if err != nil {
		log.Fatal(fmt.Errorf("redisContainer.Host %s", err))
	}

	uri := fmt.Sprintf("%s:%s", hostIP, redisContainerPort.Port())
	log.Printf("TestContainer Redis %s\n", uri)

	rdb = connect.RedisConn(ctx, configPath)
	os.Exit(m.Run())
}

func instanceReposity(t *testing.T) *Repository {
	t.Helper()
	logging := logger.NewAPILogger("debug")
	return NewRepository(rdb, dbx, logging)
}
