package command

import (
	"auth/internal/domain/entity"
	"auth/internal/infrastructure/db/gorm/dao"
	"auth/internal/infrastructure/db/gorm/mapper"
	"context"
)

type CredentialCommand struct {
	q *dao.Query
}

func NewCredentialCommand(q *dao.Query) *CredentialCommand {
	return &CredentialCommand{q: q}
}

func (c *CredentialCommand) Create(ctx context.Context, credential *entity.Credential) error {
	model := mapper.CredentialEntityToModel(credential)
	return c.q.WithContext(ctx).Credential.Create(model)
}

func (c *CredentialCommand) Update(ctx context.Context, credential *entity.Credential) error {
	model := mapper.CredentialEntityToModel(credential)
	_, err := c.q.WithContext(ctx).Credential.Where(c.q.Credential.ID.Eq(model.ID)).Updates(model)

	return err
}
