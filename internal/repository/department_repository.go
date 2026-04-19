package repository

import (
	"context"

	"github.com/auhmaugmaufm/predict-ticket-department-backend/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type departmentRepository struct {
	db *gorm.DB
}

func NewDepartmentRepository(db *gorm.DB) domain.DepartmentRepository {
	return &departmentRepository{db: db}
}

func (r *departmentRepository) Create(ctx context.Context, department *domain.Department) error {
	return r.db.WithContext(ctx).Create(department).Error
}

func (r *departmentRepository) CreateTx(tx interface{}, ctx context.Context, department *domain.Department) error {
	return tx.(*gorm.DB).WithContext(ctx).Create(department).Error
}

func (r *departmentRepository) CreateBulk(ctx context.Context, departments []domain.Department) error {
	return r.db.WithContext(ctx).CreateInBatches(departments, 100).Error
}

func (r *departmentRepository) GetByCompanyID(ctx context.Context, company_id uuid.UUID) ([]domain.Department, error) {
	var departments []domain.Department
	err := r.db.WithContext(ctx).Where("company_id = ?", company_id).Find(&departments).Error
	if err != nil {
		return nil, err
	}
	return departments, nil
}
