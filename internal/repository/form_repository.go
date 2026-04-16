package repository

import (
	"context"

	"github.com/auhmaugmaufm/predict-ticket-department-backend/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type formRepository struct {
	db *gorm.DB
}

func NewFormRepository(db *gorm.DB) domain.FormRepository {
	return &formRepository{db: db}
}

func (r *formRepository) Create(ctx context.Context, form *domain.Form) error {
	return r.db.WithContext(ctx).Create(form).Error
}

func (r *formRepository) GetByCompanyID(ctx context.Context, company_id uuid.UUID) ([]domain.Form, error) {
	var forms []domain.Form
	err := r.db.WithContext(ctx).
		Joins("Join links ON links.id = forms.form_id").
		Where("links.company_id = ?", company_id).Find(&forms).Error
	if err != nil {
		return nil, err
	}
	return forms, nil
}

func (r *formRepository) GetPerDayByCompanyID(ctx context.Context, company_id uuid.UUID, dateStr string) ([]domain.Form, error) {
	var forms []domain.Form
	err := r.db.WithContext(ctx).
		Select("forms.id", "link_id", "title", "description").
		Joins("Join links ON links.id = forms.link_id").
		Where("DATE(forms.created_at) = ?", dateStr).
		Where("links.company_id = ?", company_id).Find(&forms).Error
	if err != nil {
		return nil, err
	}
	return forms, nil
}
