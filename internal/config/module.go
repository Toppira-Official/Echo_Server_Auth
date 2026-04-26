package config

import (
	"auth/internal/application/service"
	"auth/internal/config/env"
	accesstoken "auth/internal/infrastructure/access_token"
	"auth/internal/infrastructure/cache"
	"auth/internal/infrastructure/db"

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
		func(cfg *env.Config) cache.RedisClientConfig {
			return cache.RedisClientConfig{
				Host:     cfg.Cache.RedisHost,
				Password: cfg.Cache.RedisPassword,
				Port:     cfg.Cache.RedisPort,
				DB:       cfg.Cache.RedisDB,
			}
		},
		func(cfg *env.Config) db.PostgresGormDBConfig {
			return db.PostgresGormDBConfig{
				Host:     cfg.DB.PostgresHost,
				Port:     cfg.DB.PostgresPort,
				User:     cfg.DB.PostgresUser,
				Password: cfg.DB.PostgresPassword,
				DB:       cfg.DB.PostgresDB,
			}
		},
	),
)
