package contract

import (
	"auth/internal/domain/entity"
	"auth/internal/domain/vo"
	"context"
)

type CredentialQuery interface {
	FindByID(ctx context.Context, id vo.UserID) (credential *entity.Credential, err error)
	FindByUsername(ctx context.Context, username string) (credential *entity.Credential, err error)
	ExistsByUsername(ctx context.Context, username string) (exists bool, err error)
}

type CredentialCommand interface {
	Create(ctx context.Context, credential *entity.Credential) (err error)
	Update(ctx context.Context, credential *entity.Credential) (err error)
}
