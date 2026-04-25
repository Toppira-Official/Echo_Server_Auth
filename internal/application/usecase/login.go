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
	clock               contract.Clock
}

func NewLogin(
	credentialQuery contract.CredentialQuery,
	passwordEncoder contract.PasswordEncoder,
	accessTokenSigner contract.AccessTokenSigner,
	refreshTokenFactory contract.RefreshTokenFactory,
	clock contract.Clock,
) *Login {
	return &Login{
		credentialQuery:     credentialQuery,
		passwordEncoder:     passwordEncoder,
		accessTokenSigner:   accessTokenSigner,
		refreshTokenFactory: refreshTokenFactory,
		clock:               clock,
	}
}

type LoginInput struct {
	Username string
	Password string
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

	return output, nil
}
