package usecase

import (
	"auth/internal/application/service"
	"auth/internal/domain/contract"
	"auth/internal/domain/entity"
	"auth/internal/domain/vo"
	"context"
	"errors"
)

var (
	ErrRegisterUsernameAlreadyExists = errors.New("username already exists")
)

type Register struct {
	credentialQuery   contract.CredentialQuery
	credentialCommand contract.CredentialCommand
	passwordEncoder   contract.PasswordEncoder
	uuidGenerator     contract.UuidGenerator
	clock             contract.Clock
	session           *service.Session
}

func NewRegister(
	credentialQuery contract.CredentialQuery,
	credentialCommand contract.CredentialCommand,
	passwordEncoder contract.PasswordEncoder,
	uuidGenerator contract.UuidGenerator,
	clock contract.Clock,
	session *service.Session,
) *Register {
	return &Register{
		credentialQuery:   credentialQuery,
		credentialCommand: credentialCommand,
		passwordEncoder:   passwordEncoder,
		uuidGenerator:     uuidGenerator,
		clock:             clock,
		session:           session,
	}
}

type RegisterInput struct {
	Username  string
	Password  string
	UserAgent string
	IpAddress string
}
type RegisterOutput struct {
	service.SessionTokens
}

func (r *Register) Execute(ctx context.Context, input RegisterInput) (output RegisterOutput, err error) {
	newCredential, err := r.authenticate(ctx, input.Username, input.Password)
	if err != nil {
		return output, err
	}

	tokens, err := r.session.Create(ctx, newCredential.ID(), input.UserAgent, input.IpAddress)
	if err != nil {
		return output, err
	}

	return RegisterOutput{
		SessionTokens: tokens,
	}, nil
}

func (r *Register) authenticate(ctx context.Context, username, password string) (*entity.Credential, error) {
	credential, err := r.credentialQuery.FindByUsername(ctx, username)
	if err != nil {
		if credential != nil {
			return nil, ErrRegisterUsernameAlreadyExists
		}
		return nil, err
	}

	hashedPassword, err := r.passwordEncoder.Hash(password)
	if err != nil {
		return nil, err
	}

	credentialUUID, err := r.uuidGenerator.Generate()
	if err != nil {
		return nil, err
	}

	credentialID, err := vo.NewCredentialID(credentialUUID)
	if err != nil {
		return nil, err
	}

	now := r.clock.NowUTC()

	newCredential, err := entity.NewCredential(credentialID, username, now, hashedPassword)
	if err != nil {
		return nil, err
	}

	err = r.credentialCommand.Create(ctx, newCredential)
	if err != nil {
		return nil, err
	}

	return newCredential, nil
}
