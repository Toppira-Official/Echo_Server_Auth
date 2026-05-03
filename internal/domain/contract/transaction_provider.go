package contract

import "context"

type TransactionProvider interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}
