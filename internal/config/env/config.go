package env

import (
	"os"
	"time"

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

	AccessTokenTTL  time.Duration `env:"ACCESS_TOKEN_TTL" env-default:"5h"`
	RefreshTokenTTL time.Duration `env:"REFRESH_TOKEN_TTL" env-default:"5d"`
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
