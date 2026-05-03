package transaction

import (
	"context"

	"gorm.io/gorm"
)

type GormTransaction struct {
	db *gorm.DB
}

func NewGormTransaction(db *gorm.DB) *GormTransaction {
	return &GormTransaction{db: db}
}

func (t *GormTransaction) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	if _, ok := TxFrom(ctx); ok {
		return fn(ctx)
	}

	return t.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txCtx := withTx(ctx, tx)
		return fn(txCtx)
	})
}
