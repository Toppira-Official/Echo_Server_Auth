package env

import (
	"os"
	"time"

	"github.com/Ali127Dev/xerr"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	App   AppConfig   `env-prefix:"APP_"`
	Auth  AuthConfig  `env-prefix:"AUTH_"`
	DB    DBConfig    `env-prefix:"DB_"`
	Cache CacheConfig `env-prefix:"CACHE_"`
}

type AppConfig struct {
	Port       int      `env:"PORT" env-required:"true"`
	Mode       string   `env:"MODE" env-default:"dev"`
	AppOrigins []string `env:"ORIGINS" env-separator:","`
}

type AuthConfig struct {
	TokenSecret []byte `env:"TOKEN_SECRET" env-required:"true"`

	AccessTokenTTL  time.Duration `env:"ACCESS_TOKEN_TTL" env-default:"5h"`
	RefreshTokenTTL time.Duration `env:"REFRESH_TOKEN_TTL" env-default:"120h"`
}

type DBConfig struct {
	PostgresUser     string `env:"POSTGRES_USER" env-required:"true"`
	PostgresPassword string `env:"POSTGRES_PASSWORD" env-required:"true"`
	PostgresHost     string `env:"POSTGRES_HOST" env-default:"localhost"`
	PostgresPort     int    `env:"POSTGRES_PORT" env-required:"true"`
	PostgresDB       string `env:"POSTGRES_DB" env-required:"true"`
}

type CacheConfig struct {
	RedisPassword string `env:"REDIS_PASSWORD" env-required:"true"`
	RedisHost     string `env:"REDIS_HOST" env-default:"localhost"`
	RedisPort     int    `env:"REDIS_PORT" env-required:"true"`
	RedisDB       int    `env:"REDIS_DB" env-required:"true"`
}

func NewConfig() (*Config, error) {
	if os.Getenv("APP_MODE") != "production" {
		_ = godotenv.Load()
	}

	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternalError, xerr.WithMessage("failed to initialize envs"))
	}

	return &cfg, nil
}
