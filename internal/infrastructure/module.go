package infrastructure

import (
	"auth/internal/domain/contract"
	accesstoken "auth/internal/infrastructure/access_token"
	"auth/internal/infrastructure/cache/redis"
	"auth/internal/infrastructure/clock"
	"auth/internal/infrastructure/db/gorm/command"
	"auth/internal/infrastructure/db/gorm/daoquery"
	"auth/internal/infrastructure/db/gorm/query"
	"auth/internal/infrastructure/db/postgres"
	"auth/internal/infrastructure/password"
	refreshtoken "auth/internal/infrastructure/refresh_token"
	"auth/internal/infrastructure/server/gin"
	"auth/internal/infrastructure/uuid"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"infrastructure",
	fx.Provide(
		redis.NewClient,
		postgres.NewDB,
		gin.NewGinEngine,
		daoquery.NewDaoQuery,
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
		fx.Annotate(
			redis.NewCache,
			fx.As(new(contract.Cache)),
		),
		fx.Annotate(
			command.NewCredentialCommand,
			fx.As(new(contract.CredentialCommand)),
		),
		fx.Annotate(
			query.NewCredentialQuery,
			fx.As(new(contract.CredentialQuery)),
		),
	),
)
