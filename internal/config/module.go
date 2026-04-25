package config

import (
	"auth/internal/application/service"
	"auth/internal/config/env"

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
	),
)
