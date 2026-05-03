package command

import (
	"auth/internal/domain/entity"
	"auth/internal/infrastructure/db/gorm/dao"
	"auth/internal/infrastructure/db/gorm/daoquery"
	"auth/internal/infrastructure/db/gorm/mapper"
	"context"
	"errors"

	"github.com/Ali127Dev/xerr"
	"gorm.io/gorm"
)

type CredentialCommand struct {
	q *dao.Query
}

func NewCredentialCommand(q *dao.Query) *CredentialCommand {
	return &CredentialCommand{q: q}
}

func (c *CredentialCommand) Create(ctx context.Context, credential *entity.Credential) error {
	q := daoquery.ResolveQuery(ctx, c.q)

	model := mapper.CredentialEntityToModel(credential)

	err := q.WithContext(ctx).Credential.Create(model)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return xerr.Wrap(
				err,
				xerr.CodeDuplicateKey,
				xerr.WithMessage("credential already exists"),
				xerr.WithMeta("entity", "credential"),
			)
		}

		return xerr.Wrap(
			err,
			xerr.CodeDatabaseError,
			xerr.WithMessage("failed to create credential"),
			xerr.WithMeta("entity", "credential"),
			xerr.WithMeta("operation", "create"),
		)
	}

	return nil
}

func (c *CredentialCommand) Update(ctx context.Context, credential *entity.Credential) error {
	q := daoquery.ResolveQuery(ctx, c.q)

	model := mapper.CredentialEntityToModel(credential)

	_, err := q.
		WithContext(ctx).
		Credential.
		Where(q.Credential.ID.Eq(model.ID)).
		Updates(model)

	if err != nil {
		return xerr.Wrap(
			err,
			xerr.CodeDatabaseError,
			xerr.WithMessage("failed to update credential"),
			xerr.WithMeta("entity", "credential"),
			xerr.WithMeta("operation", "update"),
			xerr.WithMeta("id", model.ID),
		)
	}

	return nil
}
