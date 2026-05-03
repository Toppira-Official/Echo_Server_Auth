package usecase

import (
	"auth/internal/application/service"
	"auth/internal/domain/contract"
	"auth/internal/domain/entity"
	"auth/internal/domain/event"
	"auth/internal/domain/vo"
	"context"
)

type Register struct {
	credentialQuery   contract.CredentialQuery
	credentialCommand contract.CredentialCommand
	passwordEncoder   contract.PasswordEncoder
	uuidGenerator     contract.UuidGenerator
	clock             contract.Clock
	session           *service.Session
	eventDispatcher   contract.EventDispatcher
	tx                contract.TransactionProvider
}

func NewRegister(
	credentialQuery contract.CredentialQuery,
	credentialCommand contract.CredentialCommand,
	passwordEncoder contract.PasswordEncoder,
	uuidGenerator contract.UuidGenerator,
	clock contract.Clock,
	session *service.Session,
	eventDispatcher contract.EventDispatcher,
	tx contract.TransactionProvider,
) *Register {
	return &Register{
		credentialQuery:   credentialQuery,
		credentialCommand: credentialCommand,
		passwordEncoder:   passwordEncoder,
		uuidGenerator:     uuidGenerator,
		clock:             clock,
		session:           session,
		eventDispatcher:   eventDispatcher,
		tx:                tx,
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
	var newCredential *entity.Credential
	var tokens service.SessionTokens

	err = r.tx.Do(ctx, func(txCtx context.Context) error {
		cred, err := r.authenticate(txCtx, input.Username, input.Password)
		if err != nil {
			return err
		}
		newCredential = cred

		tok, err := r.session.Create(txCtx, cred.ID(), input.UserAgent, input.IpAddress)
		if err != nil {
			return err
		}
		tokens = tok

		return nil
	})
	if err != nil {
		return output, err
	}

	event := event.UserRegistered{
		UserID:     newCredential.ID().Value(),
		Username:   newCredential.Username(),
		OccurredAt: r.clock.NowUTC(),
	}
	if err := r.eventDispatcher.Dispatch(ctx, event); err != nil {
		return output, err
	}

	return RegisterOutput{
		SessionTokens: tokens,
	}, nil
}

func (r *Register) authenticate(ctx context.Context, username, password string) (*entity.Credential, error) {
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
