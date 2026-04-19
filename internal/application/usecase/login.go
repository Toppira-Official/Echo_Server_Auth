package usecase

import (
	"auth/internal/domain/contract"
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
}

func NewLogin(
	credentialQuery contract.CredentialQuery,
	passwordEncoder contract.PasswordEncoder,
	accessTokenSigner contract.AccessTokenSigner,
	refreshTokenFactory contract.RefreshTokenFactory,
) *Login {
	return &Login{
		credentialQuery:     credentialQuery,
		passwordEncoder:     passwordEncoder,
		accessTokenSigner:   accessTokenSigner,
		refreshTokenFactory: refreshTokenFactory,
	}
}

type LoginInput struct {
	Username string
	Password string
}

func (l *Login) Execute(ctx context.Context, input LoginInput) (accessToken, refreshToken string, err error) {
	credential, err := l.credentialQuery.FindByUsername(ctx, input.Username)
	if err != nil {
		return "", "", ErrLoginInvalidCredentials
	}

	if ok, err := l.passwordEncoder.Compare(input.Password, credential.HashedPassword()); err != nil || !ok {
		return "", "", ErrLoginInvalidCredentials
	}

	now := time.Now().UTC()
	expiresAt := now.Add(8 * time.Hour) // TODO: must come from envs
	accessTokenPayload, err := vo.NewAccessTokenPayload(credential.ID(), now, expiresAt)
	if err != nil {
		return "", "", err
	}

	accessToken, err = l.accessTokenSigner.Generate(accessTokenPayload)
	if err != nil {
		return "", "", err
	}

	refreshToken, err = l.refreshTokenFactory.Generate()
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
