package service

import (
	"context"
	"fmt"

	"github.com/auhmaugmaufm/predict-ticket-department-backend/internal/domain"
	"github.com/google/uuid"
)

type LinkService struct {
	repo domain.LinkRepository
}

func NewLinkService(repo domain.LinkRepository) *LinkService {
	return &LinkService{repo: repo}
}

func (s *LinkService) CreateLink(ctx context.Context, company_id uuid.UUID) error {
	company_form := &domain.Link{
		CompanyID: company_id,
		Link:      GenerateLink(company_id.String()),
	}
	return s.repo.Create(ctx, company_form)
}

func (s *LinkService) GetLinkByCompanyID(ctx context.Context, company_id uuid.UUID) (*domain.Link, error) {
	company_form, err := s.repo.GetByID(ctx, company_id)
	if err != nil {
		return nil, err
	}
	return company_form, nil
}

func GenerateLink(company_id string) string {
	return fmt.Sprintf("/form/%s", company_id)
}
