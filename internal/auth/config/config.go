package config

import (
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/omekov/superapp-backend/pkg/logger"
	"gopkg.in/yaml.v3"
)

// Config ...
type Config struct {
	logg    *logger.APILogger
	Server  Server
	Logger  Logger       `yaml:"logger"`
	JWT     JWT          `yaml:"jwt"`
	GRPC    GRPC         `yaml:"grpc"`
	Mailer  MailerConfig `yaml:"mailer"`
	Migrate Migrate      `yaml:"migrate"`
}

// Logger ...
type Logger struct {
	Development       bool
	DisableCaller     bool
	DisableStacktrace bool
	Encoding          string
	Level             string `yaml:"level"`
}

// MailerConfig ...
type MailerConfig struct {
	Timeout      time.Duration `yaml:"timeout"`
	Host         string        `yaml:"host"`
	Port         int           `yaml:"port"`
	Username     string        `yaml:"username"`
	Password     string        `yaml:"password"`
	Sender       string        `yaml:"sender"`
	TemplatePath string        `yaml:"templatePath"`
}

// Server ...
type Server struct {
	Mode string
}

// Postgres ...
type Postgres struct {
	Postgres PostgresConfig `yaml:"postgres"`
}

// PostgresConfig ...
type PostgresConfig struct {
	Driver       string `yaml:"driver"`
	Host         string `yaml:"host"`
	Username     string `yaml:"username"`
	Password     string `yaml:"password"`
	DatabaseName string `yaml:"databasename"`
	Port         string `yaml:"port"`
	SSLMode      string `yaml:"sslmode"`
	MaxIdleCon   string `yaml:"maxIdleCon"`
	MaxOpenCon   string `yaml:"maxOpenCon"`
}

// Redis ...
type Redis struct {
	Redis RedisConfig `yaml:"redis"`
}

// RedisConfig ...
type RedisConfig struct {
	Host     string `yaml:"host"`
	Password string `yaml:"password"`
}

// JWT ...
type JWT struct {
	Access  string `yaml:"accent"`
	Refresh string `yaml:"refresh"`
	Mail    string `yaml:"mail"`
}

// GRPC ...
type GRPC struct {
	MaxConnectionIdle time.Duration `yaml:"maxConnectionIdle"`
	Timeout           time.Duration `yaml:"timeout"`
	MaxConnectionAge  time.Duration `yaml:"maxConnectionAge"`
	Time              time.Duration `yaml:"time"`
	TLS               bool          `yaml:"tls"`
	CertFile          string        `yaml:"certFile"`
	KeyFile           string        `yaml:"keyFile"`
}

// Migrate ...
type Migrate struct {
	Onwork   bool   `yaml:"onwork"`
	AuthPath string `yaml:"auth"`
}

// New ...
func New(logg *logger.APILogger) *Config {
	return &Config{logg: logg}
}

// GetPostgres ...
func (c *Config) GetPostgres(path string) *Postgres {
	filename, err := filepath.Abs(path)
	if err != nil {
		c.logg.Errorf("filepath.Abs: %s", err)
		return nil
	}

	configFile, err := ioutil.ReadFile(filename)
	if err != nil {
		c.logg.Errorf("ioutil.ReadFile: %s", err)
		return nil
	}

	var postgresConfig Postgres
	if err := yaml.Unmarshal(configFile, &postgresConfig); err != nil {
		c.logg.Errorf("yaml.Unmarshal(: %s", err)
		return nil
	}

	return &postgresConfig
}

// GetRedis ...
func (c *Config) GetRedis(path string) *Redis {
	filename, err := filepath.Abs(path)
	if err != nil {
		c.logg.Errorf("filepath.Abs: %s", err)
		return nil
	}

	configFile, err := ioutil.ReadFile(filename)
	if err != nil {
		c.logg.Errorf("ioutil.ReadFile: %s", err)
		return nil
	}

	var redisConfig Redis
	if err := yaml.Unmarshal(configFile, &redisConfig); err != nil {
		c.logg.Errorf("yaml.Unmarshal(: %s", err)
		return nil
	}

	return &redisConfig
}

// Get ...
func (c *Config) Get(path string) (*Config, error) {
	filename, err := filepath.Abs(path)
	if err != nil {
		c.logg.Errorf("filepath.Abs: %s", err)
		return nil, err
	}

	configFile, err := ioutil.ReadFile(filename)
	if err != nil {
		c.logg.Errorf("ioutil.ReadFile: %s", err)
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(configFile, &config); err != nil {
		c.logg.Errorf("yaml.Unmarshal(: %s", err)
		return nil, err
	}

	return &config, nil
}
