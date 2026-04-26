package config

import (
	"auth/internal/application/service"
	"auth/internal/config/env"
	accesstoken "auth/internal/infrastructure/access_token"
	"auth/internal/infrastructure/cache/redis"
	"auth/internal/infrastructure/db/postgres"
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
		func(cfg *env.Config) gin.GinEngineConfig {
			return gin.GinEngineConfig{
				Mode: cfg.App.Mode,
				Port: cfg.App.Port,
			}
		},
	),
)
