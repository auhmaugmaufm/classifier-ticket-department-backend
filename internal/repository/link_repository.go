package repository

import (
	"context"
	"errors"

	"github.com/auhmaugmaufm/predict-ticket-department-backend/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LinkRepository struct {
	db *gorm.DB
}

func NewLinkRepository(db *gorm.DB) domain.LinkRepository {
	return &LinkRepository{db: db}
}

func (r *LinkRepository) Create(ctx context.Context, Link *domain.Link) error {
	return r.db.WithContext(ctx).Create(Link).Error
}

func (r *LinkRepository) GetByID(ctx context.Context, company_id uuid.UUID) (*domain.Link, error) {
	var company_form *domain.Link
	err := r.db.WithContext(ctx).Where("company_id = ?", company_id).First(&company_form).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return company_form, nil
}
