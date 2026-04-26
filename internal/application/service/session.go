package service

import (
	"auth/internal/domain/contract"
	"auth/internal/domain/entity"
	"auth/internal/domain/vo"
	"context"
	"time"
)

type SessionConfig struct {
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

type Session struct {
	cache                 contract.Cache
	accessTokenTTL        time.Duration
	refreshTokenTTL       time.Duration
	clock                 contract.Clock
	accessTokenSigner     contract.AccessTokenSigner
	refreshTokenGenerator contract.RefreshTokenGenerator
	uuidGenerator         contract.UuidGenerator
}

func NewSession(
	cache contract.Cache,
	clock contract.Clock,
	accessTokenSigner contract.AccessTokenSigner,
	refreshTokenGenerator contract.RefreshTokenGenerator,
	uuidGenerator contract.UuidGenerator,
	cfg SessionConfig,
) *Session {
	return &Session{
		cache:                 cache,
		accessTokenTTL:        cfg.AccessTokenTTL,
		refreshTokenTTL:       cfg.RefreshTokenTTL,
		clock:                 clock,
		accessTokenSigner:     accessTokenSigner,
		refreshTokenGenerator: refreshTokenGenerator,
		uuidGenerator:         uuidGenerator,
	}
}

type SessionTokens struct {
	AccessToken           string
	RefreshToken          string
	AccessTokenExpiresAt  time.Time
	RefreshTokenExpiresAt time.Time
}

func (s *Session) Create(
	ctx context.Context, credentialID vo.CredentialID,
	userAgent, ipAddress string,
) (output SessionTokens, err error) {
	now := s.clock.NowUTC()

	refreshToken, err := s.refreshTokenGenerator.Generate()
	if err != nil {
		return output, err
	}

	deviceUUID, err := s.uuidGenerator.Generate()
	if err != nil {
		return output, err
	}

	deviceID, err := vo.NewDeviceID(deviceUUID)
	if err != nil {
		return output, err
	}

	refreshTokenExpiresAt := now.Add(s.refreshTokenTTL)

	newDevice, err := entity.NewDevice(
		deviceID, refreshToken, refreshTokenExpiresAt,
		now, userAgent, ipAddress,
	)
	if err != nil {
		return output, err
	}

	cacheKey := "refresh:" + refreshToken
	if err := s.cache.Set(ctx, cacheKey, newDevice, s.refreshTokenTTL); err != nil {
		return output, err
	}

	accessTokenExpiresAt := now.Add(s.accessTokenTTL)
	accessTokenPayload, err := vo.NewAccessTokenPayload(credentialID, now, accessTokenExpiresAt)
	if err != nil {
		return output, err
	}

	accessToken, err := s.accessTokenSigner.Generate(accessTokenPayload)
	if err != nil {
		return output, err
	}

	return SessionTokens{
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessTokenExpiresAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshTokenExpiresAt,
	}, nil
}
