package repository

import (
	"context"
	"errors"

	"github.com/auhmaugmaufm/predict-ticket-department-backend/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type companyFormRepository struct {
	db *gorm.DB
}

func NewCompanyFormRepository(db *gorm.DB) domain.CompanyFormRepository {
	return &companyFormRepository{db: db}
}

func (r *companyFormRepository) Create(ctx context.Context, companyForm *domain.CompanyForm) error {
	return r.db.WithContext(ctx).Create(companyForm).Error
}

func (r *companyFormRepository) GetByID(ctx context.Context, company_id uuid.UUID) (*domain.CompanyForm, error) {
	var company_form *domain.CompanyForm
	err := r.db.WithContext(ctx).Where("company_id = ?", company_id).First(&company_form).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return company_form, nil
}
