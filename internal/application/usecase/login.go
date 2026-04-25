package usecase

import (
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
}

func NewLogin(
	credentialQuery contract.CredentialQuery,
	passwordEncoder contract.PasswordEncoder,
	accessTokenSigner contract.AccessTokenSigner,
	refreshTokenFactory contract.RefreshTokenFactory,
	uuidGenerator contract.UuidGenerator,
	clock contract.Clock,
	cache contract.Cache,
) *Login {
	return &Login{
		credentialQuery:     credentialQuery,
		passwordEncoder:     passwordEncoder,
		accessTokenSigner:   accessTokenSigner,
		refreshTokenFactory: refreshTokenFactory,
		uuidGenerator:       refreshTokenFactory,
		clock:               clock,
		cache:               cache,
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

	now := l.clock.Now_UTC()
	expiresAt := now.Add(8 * time.Hour) // TODO: must come from envs
	accessTokenPayload, err := vo.NewAccessTokenPayload(credential.ID(), now, expiresAt)
	if err != nil {
		return output, err
	}

	accessToken, err := l.accessTokenSigner.Generate(accessTokenPayload)
	if err != nil {
		return output, err
	}

	refreshToken, err := l.refreshTokenFactory.Generate()
	if err != nil {
		return output, err
	}

	output.AccessToken = accessToken
	output.RefreshToken = refreshToken

	deviceUUID, err := l.uuidGenerator.Generate()
	if err != nil {
		return output, err
	}

	deviceID, err := vo.NewDeviceID(deviceUUID)
	if err != nil {
		return output, err
	}

	expiresAt = now.Add(8 * time.Hour) // TODO: must come from envs

	newDevice, err := entity.NewDevice(
		deviceID, refreshToken, expiresAt,
		now, input.UserAgent, input.IpAddress,
	)
	if err != nil {
		return output, err
	}

	cacheKey := "device:" + input.UserAgent
	if err := l.cache.Set(ctx, cacheKey, newDevice, expiresAt.Sub(now)); err != nil {
		return output, err
	}

	return output, nil
}
