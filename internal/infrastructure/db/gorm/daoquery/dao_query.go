package daoquery

import (
	"auth/internal/infrastructure/db/gorm/dao"
	"auth/internal/infrastructure/db/gorm/transaction"
	"context"

	"gorm.io/gorm"
)

func NewDaoQuery(db *gorm.DB) *dao.Query {
	return dao.Use(db)
}

func ResolveQuery(ctx context.Context, q *dao.Query) *dao.Query {
	if tx, ok := transaction.TxFrom(ctx); ok {
		return dao.Use(tx)
	}
	return q
}
