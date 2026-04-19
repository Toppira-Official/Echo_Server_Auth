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
	ErrRegisterUsernameAlreadyExists = errors.New("username already exists")
)

type Register struct {
	credentialQuery     contract.CredentialQuery
	credentialCommand   contract.CredentialCommand
	passwordEncoder     contract.PasswordEncoder
	accessTokenSigner   contract.AccessTokenSigner
	refreshTokenFactory contract.RefreshTokenFactory
	uuidGenerator       contract.UuidGenerator
}

func NewRegister(
	credentialQuery contract.CredentialQuery,
	credentialCommand contract.CredentialCommand,
	passwordEncoder contract.PasswordEncoder,
	accessTokenSigner contract.AccessTokenSigner,
	refreshTokenFactory contract.RefreshTokenFactory,
	uuidGenerator contract.UuidGenerator,
) *Register {
	return &Register{
		credentialQuery:     credentialQuery,
		credentialCommand:   credentialCommand,
		passwordEncoder:     passwordEncoder,
		accessTokenSigner:   accessTokenSigner,
		refreshTokenFactory: refreshTokenFactory,
		uuidGenerator:       uuidGenerator,
	}
}

type RegisterInput struct {
	Username string
	Password string
}

func (r *Register) Execute(ctx context.Context, input RegisterInput) (accessToken, refreshToken string, err error) {
	credential, err := r.credentialQuery.FindByUsername(ctx, input.Username)
	if err != nil {
		if credential != nil {
			return "", "", ErrRegisterUsernameAlreadyExists
		}
		return "", "", err
	}

	hashedPassword, err := r.passwordEncoder.Hash(input.Password)
	if err != nil {
		return "", "", err
	}

	uuid, err := r.uuidGenerator.Generate()
	if err != nil {
		return "", "", err
	}

	credentialID, err := vo.NewCredentialID(uuid)
	if err != nil {
		return "", "", err
	}

	now := time.Now().UTC()

	newCredential, err := entity.NewCredential(credentialID, input.Username, now, hashedPassword)
	if err != nil {
		return "", "", err
	}

	err = r.credentialCommand.Create(ctx, newCredential)
	if err != nil {
		return "", "", err
	}

	expiresAt := now.Add(8 * time.Hour) // TODO: must come from envs
	accessTokenPayload, err := vo.NewAccessTokenPayload(newCredential.ID(), now, expiresAt)
	if err != nil {
		return "", "", err
	}

	accessToken, err = r.accessTokenSigner.Generate(accessTokenPayload)
	if err != nil {
		return "", "", err
	}

	refreshToken, err = r.refreshTokenFactory.Generate()
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
