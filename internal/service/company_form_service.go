package service

import (
	"context"
	"fmt"

	"github.com/auhmaugmaufm/predict-ticket-department-backend/internal/domain"
	"github.com/google/uuid"
)

type CompanyFormService struct {
	repo domain.CompanyFormRepository
}

func NewCompanyFormService(repo domain.CompanyFormRepository) *CompanyFormService {
	return &CompanyFormService{repo: repo}
}

func (s *CompanyFormService) CreateCompanyForm(ctx context.Context, company_id uuid.UUID) error {
	company_form := &domain.CompanyForm{
		CompanyID: company_id,
		LinkForm:  GenerateLink(company_id.String()),
	}
	return s.repo.Create(ctx, company_form)
}

func (s *CompanyFormService) GetCompanyFormByCompanyID(ctx context.Context, company_id uuid.UUID) (*domain.CompanyForm, error) {
	company_form, err := s.repo.GetByID(ctx, company_id)
	if err != nil {
		return nil, err
	}
	return company_form, nil
}

func GenerateLink(company_id string) string {
	return fmt.Sprintf("/form/%s", company_id)
}
