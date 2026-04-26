package command

import (
	"auth/internal/domain/entity"
	"auth/internal/infrastructure/db/gorm/dao"
	"auth/internal/infrastructure/db/gorm/mapper"
	"context"
	"errors"

	"github.com/Ali127Dev/xerr"
	"github.com/jackc/pgx/v5/pgconn"
)

type CredentialCommand struct {
	q *dao.Query
}

func NewCredentialCommand(q *dao.Query) *CredentialCommand {
	return &CredentialCommand{q: q}
}

func (c *CredentialCommand) Create(ctx context.Context, credential *entity.Credential) error {
	model := mapper.CredentialEntityToModel(credential)

	err := c.q.WithContext(ctx).Credential.Create(model)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505": // unique_violation
				return xerr.Wrap(
					err,
					xerr.CodeDuplicateKey,
					xerr.WithMessage("credential already exists"),
					xerr.WithMeta("entity", "credential"),
					xerr.WithMeta("constraint", pgErr.ConstraintName),
				)
			case "23503": // foreign_key_violation
				return xerr.Wrap(
					err,
					xerr.CodeForeignKeyError,
					xerr.WithMessage("invalid foreign key reference"),
					xerr.WithMeta("entity", "credential"),
					xerr.WithMeta("constraint", pgErr.ConstraintName),
				)
			}
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
	model := mapper.CredentialEntityToModel(credential)

	_, err := c.q.
		WithContext(ctx).
		Credential.
		Where(c.q.Credential.ID.Eq(model.ID)).
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
