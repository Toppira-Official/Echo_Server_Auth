package env

import (
	"os"
	"time"

	"github.com/Ali127Dev/xerr"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	App    AppConfig    `env-prefix:"APP_"`
	Auth   AuthConfig   `env-prefix:"AUTH_"`
	DB     DBConfig     `env-prefix:"DB_"`
	Cache  CacheConfig  `env-prefix:"CACHE_"`
	Logger LoggerConfig `env-prefix:"LOG_"`
	Kafka  KafkaConfig  `env-prefix:"KAFKA_"`
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

type LoggerConfig struct {
	Level      string `env:"LEVEL" env-default:"info"`
	Path       string `env:"PATH" env-default:"./logs/app.log"`
	MaxSize    int    `env:"MAX_SIZE_MB" env-default:"50"`
	MaxBackups int    `env:"MAX_BACKUPS" env-default:"10"`
	MaxAge     int    `env:"MAX_AGE_DAYS" env-default:"30"`
	Compress   bool   `env:"COMPRESS" env-default:"true"`
}

type KafkaConfig struct {
	Brokers []string `env:"BROKERS" env-separator:","`
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
