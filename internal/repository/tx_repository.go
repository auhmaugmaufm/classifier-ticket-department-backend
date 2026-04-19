package repository

import (
	"context"

	"github.com/auhmaugmaufm/predict-ticket-department-backend/internal/domain"
	"gorm.io/gorm"
)

type txRepository struct {
	db *gorm.DB
}

func NewTxManager(db *gorm.DB) domain.TxRepository {
	return &txRepository{db: db}
}

func (t *txRepository) WithinTransaction(ctx context.Context, fn func(txCtx context.Context, tx interface{}) error) error {
	return t.db.WithContext(ctx).Transaction(func(txDB *gorm.DB) error {
		return fn(ctx, txDB)
	})
}
