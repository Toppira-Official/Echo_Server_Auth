package infrastructure

import (
	"auth/internal/domain/contract"
	accesstoken "auth/internal/infrastructure/access_token"
	"auth/internal/infrastructure/cache"
	"auth/internal/infrastructure/clock"
	"auth/internal/infrastructure/password"
	refreshtoken "auth/internal/infrastructure/refresh_token"
	"auth/internal/infrastructure/uuid"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"infrastructure",
	fx.Provide(
		cache.NewRedisClient,
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
			fx.As(new(contract.RefreshTokenGenerator)),
		),
		fx.Annotate(
			refreshtoken.NewSha256RefreshTokenHasher,
			fx.As(new(contract.RefreshTokenHasher)),
		),
		fx.Annotate(
			accesstoken.NewJwtAccessTokenSigner,
			fx.As(new(contract.AccessTokenSigner)),
		),
	),
)
