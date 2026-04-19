package domain

import "context"

type TxRepository interface {
	WithinTransaction(ctx context.Context, fn func(txCtx context.Context, tx interface{}) error) error
}
