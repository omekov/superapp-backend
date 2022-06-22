package repository_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/docker/go-connections/nat"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	. "github.com/omekov/superapp-backend/internal/salecar/car/repository"
	"github.com/omekov/superapp-backend/pkg/logger"
	"github.com/pressly/goose/v3"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var dbx *sqlx.DB

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
			log.Fatal(fmt.Errorf("pgContainer.Terminate %s", err))
		}
	}()

	pgContainerPort, err := pgContainer.MappedPort(ctx, "5432")
	if err != nil {
		log.Fatal(fmt.Errorf("pgContainer.MappedPort %s", err))
	}

	pgContainerHostIP, err := pgContainer.Host(ctx)
	if err != nil {
		log.Fatal(fmt.Errorf("pgContainer.Host %s", err))
	}

	log.Println("TestContainer Postgres PORT:", pgContainerPort.Port())
	logg := logger.NewAPILogger("debug")
	logg.InitLogger()
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

	err = goose.Up(dbx.DB, "../../../../migrations/salecar")
	if err != nil {
		log.Fatal(fmt.Errorf("goose.Up %s", err))
	}
	os.Exit(m.Run())
}

func instanceReposity(t *testing.T) *CarRepository {
	t.Helper()
	if dbx == nil {
		panic("dbx is nil")
	}
	return NewCarRepository(dbx)
}
