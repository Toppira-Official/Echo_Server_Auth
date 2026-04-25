package usecase

import (
	"auth/internal/config/env"
	"auth/internal/domain/contract"
	"auth/internal/domain/entity"
	"auth/internal/domain/vo"
	"context"
	"errors"
	"time"
)

var (
	ErrLoginInvalidCredentials = errors.New("username or password is invalid")
)

type Login struct {
	credentialQuery     contract.CredentialQuery
	passwordEncoder     contract.PasswordEncoder
	accessTokenSigner   contract.AccessTokenSigner
	refreshTokenFactory contract.RefreshTokenFactory
	uuidGenerator       contract.UuidGenerator
	clock               contract.Clock
	cache               contract.Cache
	accessTokenTTL      time.Duration
	refreshTokenTTL     time.Duration
}

func NewLogin(
	credentialQuery contract.CredentialQuery,
	passwordEncoder contract.PasswordEncoder,
	accessTokenSigner contract.AccessTokenSigner,
	refreshTokenFactory contract.RefreshTokenFactory,
	uuidGenerator contract.UuidGenerator,
	clock contract.Clock,
	cache contract.Cache,
	envConfig env.Config,
) *Login {
	accessTokenTTL := time.Duration(envConfig.Auth.AccessTokenExpiresInHours) * time.Hour
	refreshTokenTTL := time.Duration(envConfig.Auth.RefreshTokenExpiresInDays) * 24 * time.Hour

	return &Login{
		credentialQuery:     credentialQuery,
		passwordEncoder:     passwordEncoder,
		accessTokenSigner:   accessTokenSigner,
		refreshTokenFactory: refreshTokenFactory,
		uuidGenerator:       uuidGenerator,
		clock:               clock,
		cache:               cache,
		accessTokenTTL:      accessTokenTTL,
		refreshTokenTTL:     refreshTokenTTL,
	}
}

type LoginInput struct {
	Username  string
	Password  string
	UserAgent string
	IpAddress string
}
type LoginOutput struct {
	AccessToken  string
	RefreshToken string
}

func (l *Login) Execute(ctx context.Context, input LoginInput) (output LoginOutput, err error) {
	credential, err := l.credentialQuery.FindByUsername(ctx, input.Username)
	if err != nil {
		return output, ErrLoginInvalidCredentials
	}

	if ok, err := l.passwordEncoder.Compare(input.Password, credential.HashedPassword()); err != nil || !ok {
		return output, ErrLoginInvalidCredentials
	}

	now := l.clock.NowUTC()

	refreshToken, err := l.refreshTokenFactory.Generate()
	if err != nil {
		return output, err
	}

	deviceUUID, err := l.uuidGenerator.Generate()
	if err != nil {
		return output, err
	}

	deviceID, err := vo.NewDeviceID(deviceUUID)
	if err != nil {
		return output, err
	}

	refreshTokenExpiresAt := now.Add(l.refreshTokenTTL)

	newDevice, err := entity.NewDevice(
		deviceID, refreshToken, refreshTokenExpiresAt,
		now, input.UserAgent, input.IpAddress,
	)
	if err != nil {
		return output, err
	}

	cacheKey := "refresh:" + refreshToken
	if err := l.cache.Set(ctx, cacheKey, newDevice, l.refreshTokenTTL); err != nil {
		return output, err
	}

	accessTokenExpiresAt := now.Add(l.accessTokenTTL)
	accessTokenPayload, err := vo.NewAccessTokenPayload(credential.ID(), now, accessTokenExpiresAt)
	if err != nil {
		return output, err
	}

	accessToken, err := l.accessTokenSigner.Generate(accessTokenPayload)
	if err != nil {
		return output, err
	}

	return LoginOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
