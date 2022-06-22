package config

import (
	"errors"
	"log"

	"github.com/spf13/viper"
)

// PostgresConfig ...
type PostgresConfig struct {
	PostgresqlHost     string
	PostgresqlPort     string
	PostgresqlUser     string
	PostgresqlPassword string
	PostgresqlDbname   string
	PostgresqlSSLMode  bool
	PgDriver           string
}

// Config ...
type Config struct {
	Postgres PostgresConfig
}

func loadConfig(filename string) (*viper.Viper, error) {
	v := viper.New()

	v.SetConfigName(filename)
	v.SetConfigType("yaml")
	v.AddConfigPath("./configs")
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found")
		}
		return nil, err
	}

	return v, nil
}

// Parse config file
func parseConfig(v *viper.Viper) (*Config, error) {
	var c Config

	err := v.Unmarshal(&c)
	if err != nil {
		log.Printf("unable to decode into struct, %v", err)
		return nil, err
	}

	return &c, nil
}

// Get config
func Get(configPath string) (*Config, error) {
	cfgFile, err := loadConfig(configPath)
	if err != nil {
		return nil, err
	}

	cfg, err := parseConfig(cfgFile)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
