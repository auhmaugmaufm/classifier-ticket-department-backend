package service

import (
	"context"

	"github.com/auhmaugmaufm/predict-ticket-department-backend/internal/domain"
	"github.com/google/uuid"
)

type FormService struct {
	repo domain.FormRepository
}

func NewFormService(repo domain.FormRepository) *FormService {
	return &FormService{repo: repo}
}

func (s *FormService) SubmitForm(ctx context.Context, form *domain.Form) error {
	return s.repo.Create(ctx, form)
}

func (s *FormService) GetSubmitFormByCompanyID(ctx context.Context, company_id uuid.UUID) ([]domain.Form, error) {
	forms, err := s.repo.GetByCompanyID(ctx, company_id)
	if err != nil {
		return nil, err
	}
	return forms, nil
}

func (s *FormService) GetSubmitFormPerDayByCompanyID(ctx context.Context, company_id uuid.UUID, dateStr string) ([]domain.Form, error) {
	forms, err := s.repo.GetPerDayByCompanyID(ctx, company_id, dateStr)
	if err != nil {
		return nil, err
	}
	return forms, nil
}
