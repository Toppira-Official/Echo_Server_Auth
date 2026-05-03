package usecase

import (
	"auth/internal/application/service"
	"auth/internal/domain/contract"
	"auth/internal/domain/entity"
	"context"

	"github.com/Ali127Dev/xerr"
)

var (
	ErrLoginInvalidCredentials = xerr.New(xerr.CodePermissionDenied, xerr.WithMessage("username or password is invalid"))
)

type Login struct {
	credentialQuery contract.CredentialQuery
	passwordEncoder contract.PasswordEncoder
	session         *service.Session
	tx              contract.TransactionProvider
}

func NewLogin(
	credentialQuery contract.CredentialQuery,
	passwordEncoder contract.PasswordEncoder,
	session *service.Session,
	tx contract.TransactionProvider,
) *Login {
	return &Login{
		credentialQuery: credentialQuery,
		passwordEncoder: passwordEncoder,
		session:         session,
		tx:              tx,
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
	var credential *entity.Credential
	var tokens service.SessionTokens

	err = l.tx.Do(ctx, func(ctx context.Context) error {
		cred, err := l.authenticate(ctx, input.Username, input.Password)
		if err != nil {
			return ErrLoginInvalidCredentials
		}
		credential = cred

		tok, err := l.session.Create(ctx, credential.ID(), input.UserAgent, input.IpAddress)
		if err != nil {
			return err
		}
		tokens = tok

		return nil
	})
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
