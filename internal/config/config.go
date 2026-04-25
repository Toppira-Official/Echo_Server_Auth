package config

import (
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	App  AppConfig  `env-prefix:"APP_"`
	Auth AuthConfig `env-prefix:"AUTH_"`
}

type AppConfig struct {
	Port int    `env:"PORT" env-required:"true"`
	Mode string `env:"MODE" env-default:"dev"`
}

type AuthConfig struct {
	TokenSecret string `env:"TOKEN_SECRET" env-required:"true"`

	AccessTokenExpiresInHours int `env:"ACCESS_TOKEN_EXPIRES_IN_HOURS" env-default:"5"`
	RefreshTokenExpiresInDays int `env:"REFRESH_TOKEN_EXPIRES_IN_DAYS" env-default:"5"`
}

func NewConfig() (*Config, error) {
	if os.Getenv("APP_MODE") != "production" {
		_ = godotenv.Load()
	}

	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
