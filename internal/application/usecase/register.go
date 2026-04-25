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
	ErrRegisterUsernameAlreadyExists = errors.New("username already exists")
)

type Register struct {
	credentialQuery     contract.CredentialQuery
	credentialCommand   contract.CredentialCommand
	passwordEncoder     contract.PasswordEncoder
	accessTokenSigner   contract.AccessTokenSigner
	refreshTokenFactory contract.RefreshTokenFactory
	uuidGenerator       contract.UuidGenerator
	clock               contract.Clock
	cache               contract.Cache
	accessTokenTTL      time.Duration
	refreshTokenTTL     time.Duration
}

func NewRegister(
	credentialQuery contract.CredentialQuery,
	credentialCommand contract.CredentialCommand,
	passwordEncoder contract.PasswordEncoder,
	accessTokenSigner contract.AccessTokenSigner,
	refreshTokenFactory contract.RefreshTokenFactory,
	uuidGenerator contract.UuidGenerator,
	clock contract.Clock,
	cache contract.Cache,
	envConfig env.Config,
) *Register {
	accessTokenTTL := time.Duration(envConfig.Auth.AccessTokenExpiresInHours) * time.Hour
	refreshTokenTTL := time.Duration(envConfig.Auth.RefreshTokenExpiresInDays) * 24 * time.Hour

	return &Register{
		credentialQuery:     credentialQuery,
		credentialCommand:   credentialCommand,
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

type RegisterInput struct {
	Username  string
	Password  string
	UserAgent string
	IpAddress string
}
type RegisterOutput struct {
	AccessToken  string
	RefreshToken string
}

func (r *Register) Execute(ctx context.Context, input RegisterInput) (output RegisterOutput, err error) {
	credential, err := r.credentialQuery.FindByUsername(ctx, input.Username)
	if err != nil {
		if credential != nil {
			return output, ErrRegisterUsernameAlreadyExists
		}
		return output, err
	}

	hashedPassword, err := r.passwordEncoder.Hash(input.Password)
	if err != nil {
		return output, err
	}

	credentialUUID, err := r.uuidGenerator.Generate()
	if err != nil {
		return output, err
	}

	credentialID, err := vo.NewCredentialID(credentialUUID)
	if err != nil {
		return output, err
	}

	now := r.clock.NowUTC()

	newCredential, err := entity.NewCredential(credentialID, input.Username, now, hashedPassword)
	if err != nil {
		return output, err
	}

	err = r.credentialCommand.Create(ctx, newCredential)
	if err != nil {
		return output, err
	}

	refreshToken, err := r.refreshTokenFactory.Generate()
	if err != nil {
		return output, err
	}

	deviceUUID, err := r.uuidGenerator.Generate()
	if err != nil {
		return output, err
	}

	deviceID, err := vo.NewDeviceID(deviceUUID)
	if err != nil {
		return output, err
	}

	refreshTokenExpiresAt := now.Add(r.refreshTokenTTL)

	newDevice, err := entity.NewDevice(
		deviceID, refreshToken, refreshTokenExpiresAt,
		now, input.UserAgent, input.IpAddress,
	)
	if err != nil {
		return output, err
	}

	cacheKey := "refresh:" + refreshToken
	if err := r.cache.Set(ctx, cacheKey, newDevice, r.refreshTokenTTL); err != nil {
		return output, err
	}

	accessTokenExpiresAt := now.Add(r.accessTokenTTL)
	accessTokenPayload, err := vo.NewAccessTokenPayload(newCredential.ID(), now, accessTokenExpiresAt)
	if err != nil {
		return output, err
	}

	accessToken, err := r.accessTokenSigner.Generate(accessTokenPayload)
	if err != nil {
		return output, err
	}

	return RegisterOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
