package config

import "github.com/caarlos0/env/v6"

type Config struct {
	ServerAddr          string `env:"HTTP_PORT" envDefault:":8080"`
	TelegramApiToken    string `env:"TELEGRAM_API_TOKEN,required"`
	KGDURL              string `env:"KGD_URL,required"`
	OpenExchangeRateURL string `env:"OPEN_EXCHANGE_RATE_URL,required"`
	AssetsDir           string `env:"FRONT_FILES_PATH,required"`
	SqlitePath          string `env:"SQLITE_PATH,required"`
}

func Get() (Config, error) {
	cfg := Config{}
	if err := cfg.readFromEnvironment(); err != nil {
		return cfg, err
	}
	return cfg, nil
}

// readFromEnvironment reads the settings from environment variables.
func (c *Config) readFromEnvironment() error {
	return env.Parse(c)
}
