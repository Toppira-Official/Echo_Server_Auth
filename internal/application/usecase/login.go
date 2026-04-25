package usecase

import (
	"auth/internal/application/service"
	"auth/internal/domain/contract"
	"auth/internal/domain/entity"
	"context"
	"errors"
)

var (
	ErrLoginInvalidCredentials = errors.New("username or password is invalid")
)

type Login struct {
	credentialQuery contract.CredentialQuery
	passwordEncoder contract.PasswordEncoder
	session         *service.Session
}

func NewLogin(
	credentialQuery contract.CredentialQuery,
	passwordEncoder contract.PasswordEncoder,
	session *service.Session,
) *Login {
	return &Login{
		credentialQuery: credentialQuery,
		passwordEncoder: passwordEncoder,
		session:         session,
	}
}

type LoginInput struct {
	Username  string
	Password  string
	UserAgent string
	IpAddress string
}
type LoginOutput struct {
	service.SessionTokens
}

func (l *Login) Execute(ctx context.Context, input LoginInput) (output LoginOutput, err error) {
	credential, err := l.authenticate(ctx, input.Username, input.Password)
	if err != nil {
		return output, ErrLoginInvalidCredentials
	}

	tokens, err := l.session.Create(ctx, credential.ID(), input.UserAgent, input.IpAddress)
	if err != nil {
		return output, err
	}

	return LoginOutput{
		SessionTokens: tokens,
	}, nil
}

func (l *Login) authenticate(ctx context.Context, username, password string) (*entity.Credential, error) {
	credential, err := l.credentialQuery.FindByUsername(ctx, username)
	if err != nil {
		return nil, ErrLoginInvalidCredentials
	}

	if ok, err := l.passwordEncoder.Compare(password, credential.HashedPassword()); err != nil || !ok {
		return nil, ErrLoginInvalidCredentials
	}

	return credential, nil
}
