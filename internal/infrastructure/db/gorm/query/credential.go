package query

import (
	"auth/internal/domain/entity"
	"auth/internal/domain/vo"
	"auth/internal/infrastructure/db/gorm/dao"
	"auth/internal/infrastructure/db/gorm/mapper"
	"context"
	"errors"

	"gorm.io/gorm"
)

type CredentialQuery struct {
	q *dao.Query
}

func NewCredentialQuery(q *dao.Query) *CredentialQuery {
	return &CredentialQuery{q: q}
}

func (c *CredentialQuery) FindByID(ctx context.Context, id vo.CredentialID) (*entity.Credential, error) {
	model, err := c.q.WithContext(ctx).Credential.Where(c.q.Credential.ID.Eq(id.Value())).First()
	if err != nil {
		return nil, err
	}

	return mapper.CredentialModelToEntity(model)
}

func (c *CredentialQuery) FindByUsername(ctx context.Context, username string) (credential *entity.Credential, err error) {
	model, err := c.q.WithContext(ctx).Credential.Where(c.q.Credential.Username.Eq(username)).First()
	if err != nil {
		return nil, err
	}

	return mapper.CredentialModelToEntity(model)
}

func (c *CredentialQuery) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	_, err := c.q.WithContext(ctx).Credential.
		Select(c.q.Credential.ID).
		Where(c.q.Credential.Username.Eq(username)).
		First()

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
