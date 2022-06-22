package conn

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/omekov/superapp-backend/internal/auth/config"
	"github.com/omekov/superapp-backend/pkg/logger"
)

// Conn ...
type Conn struct {
	logg             *logger.APILogger
	pgTimeAttempt    time.Duration
	redisTimeAttempt time.Duration
}

// New ...
func New(logg *logger.APILogger) *Conn {
	return &Conn{logg: logg, pgTimeAttempt: 10, redisTimeAttempt: 10}
}

// SQLXConn ...
func (c *Conn) SQLXConn(ctx context.Context, configYamlPath string) *sqlx.DB {
	dbx, err := c.sqlxInit(ctx, configYamlPath)
	if err != nil {
		c.logg.Errorf("sqlx.Connect: %s", err)
		ticker := time.NewTicker(c.pgTimeAttempt * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			dbx, err = c.sqlxInit(ctx, configYamlPath)
			if err == nil {
				break
			}
			c.logg.Errorf("Trying to connect every %d seconds; err: %s", c.pgTimeAttempt, err)
			dbx = nil
		}
	}
	return dbx
}

// RedisConn ...
func (c *Conn) RedisConn(ctx context.Context, configYamlPath string) *redis.Client {
	client, err := c.redisInit(ctx, configYamlPath)
	if err != nil {
		c.logg.Errorf("sqlx.Connect: %s", err)
		ticker := time.NewTicker(c.pgTimeAttempt * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			client, err = c.redisInit(ctx, configYamlPath)
			if err == nil {
				break
			}
			c.logg.Errorf("Trying to connect every %d seconds; err: %s", c.pgTimeAttempt, err)
			client = nil
		}
	}
	return client
}

func (c *Conn) redisInit(ctx context.Context, path string) (*redis.Client, error) {
	cfg := config.New(c.logg).GetRedis(path)
	opt := redis.Options{
		Addr:         cfg.Redis.Host,
		MinIdleConns: 10,
		PoolSize:     0,
		PoolTimeout:  time.Duration(10) * time.Second,
		Password:     cfg.Redis.Password, // no password set
		DB:           1,                  // use default DB
	}
	client := redis.NewClient(&opt)
	stmt := client.Ping(ctx)
	if err := stmt.Err(); stmt.Err() != nil {
		return nil, err
	}
	return client, nil
}

func (c *Conn) sqlxInit(ctx context.Context, path string) (*sqlx.DB, error) {
	cfg := config.New(c.logg).GetPostgres(path)
	dsn := fmt.Sprintf(
		"host=%s port=%s dbname=%s sslmode=%s user=%s password=%s",
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.DatabaseName,
		cfg.Postgres.SSLMode,
		cfg.Postgres.Username,
		cfg.Postgres.Password,
	)
	return sqlx.ConnectContext(ctx, cfg.Postgres.Driver, dsn)
}
