package service

import (
	"context"

	"github.com/auhmaugmaufm/predict-ticket-department-backend/internal/domain"
	"github.com/google/uuid"
)

type DepartmentService struct {
	repo domain.DepartmentRepository
}

func NewDepartmentService(repo domain.DepartmentRepository) *DepartmentService {
	return &DepartmentService{repo: repo}
}

func (s *DepartmentService) AddDepartments(ctx context.Context, departments []domain.Department) error {
	return s.repo.CreateBulk(ctx, departments)
}

func (s *DepartmentService) GetDepartmentsByCompanyID(ctx context.Context, company_id uuid.UUID) ([]domain.Department, error) {
	return s.repo.GetByCompanyID(ctx, company_id)
}

func (s *DepartmentService) UpdateDepartmentStatus(ctx context.Context, id uuid.UUID, isActive bool) error {
	return s.repo.UpdateStatus(ctx, id, isActive)
}
