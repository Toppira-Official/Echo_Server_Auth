package config

import (
	"auth/internal/application/service"
	"auth/internal/config/env"
	accesstoken "auth/internal/infrastructure/access_token"
	"auth/internal/infrastructure/cache/redis"
	"auth/internal/infrastructure/db/postgres"
	"auth/internal/infrastructure/kafka"
	"auth/internal/infrastructure/logger"
	"auth/internal/infrastructure/server/gin"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"config",
	fx.Provide(
		env.NewConfig,
		func(cfg *env.Config) service.SessionConfig {
			return service.SessionConfig{
				AccessTokenTTL:  cfg.Auth.AccessTokenTTL,
				RefreshTokenTTL: cfg.Auth.RefreshTokenTTL,
			}
		},
		func(cfg *env.Config) accesstoken.JwtAccessTokenSignerConfig {
			return accesstoken.JwtAccessTokenSignerConfig{
				SecretKey: cfg.Auth.TokenSecret,
			}
		},
		func(cfg *env.Config) redis.ClientConfig {
			return redis.ClientConfig{
				Host:     cfg.Cache.RedisHost,
				Password: cfg.Cache.RedisPassword,
				Port:     cfg.Cache.RedisPort,
				DB:       cfg.Cache.RedisDB,
			}
		},
		func(cfg *env.Config) postgres.DBConfig {
			return postgres.DBConfig{
				Host:     cfg.DB.PostgresHost,
				Port:     cfg.DB.PostgresPort,
				User:     cfg.DB.PostgresUser,
				Password: cfg.DB.PostgresPassword,
				DB:       cfg.DB.PostgresDB,
			}
		},
		func(cfg *env.Config) gin.EngineConfig {
			return gin.EngineConfig{
				Mode:    cfg.App.Mode,
				Port:    cfg.App.Port,
				Origins: cfg.App.AppOrigins,
			}
		},
		func(cfg *env.Config) logger.ZapLoggerConfig {
			return logger.ZapLoggerConfig{
				Mode:       cfg.App.Mode,
				LogPath:    cfg.Logger.Path,
				MaxSize:    cfg.Logger.MaxSize,
				MaxBackups: cfg.Logger.MaxBackups,
				MaxAge:     cfg.Logger.MaxAge,
				Compress:   cfg.Logger.Compress,
			}
		},
		func(cfg *env.Config) kafka.ProducerConfig {
			return kafka.ProducerConfig{
				Brokers: cfg.Kafka.Brokers,
			}
		},
	),
)
