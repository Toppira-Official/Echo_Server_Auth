package infrastructure

import (
	"auth/internal/domain/contract"
	"auth/internal/infrastructure/clock"
	"auth/internal/infrastructure/password"
	refreshtoken "auth/internal/infrastructure/refresh_token"
	"auth/internal/infrastructure/uuid"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"infrastructure",
	fx.Provide(
		fx.Annotate(
			uuid.NewKsuidIdGenerator,
			fx.As(new(contract.UuidGenerator)),
		),
		fx.Annotate(
			clock.NewSystemClock,
			fx.As(new(contract.Clock)),
		),
		fx.Annotate(
			password.NewBcryptPasswordEncoder,
			fx.As(new(contract.PasswordEncoder)),
		),
		fx.Annotate(
			refreshtoken.NewRandomRefreshTokenFactory,
			fx.As(new(contract.RefreshTokenFactory)),
		),
	),
)
