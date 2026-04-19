package repository

import (
	"context"
	"errors"

	"github.com/auhmaugmaufm/predict-ticket-department-backend/internal/domain"
	"gorm.io/gorm"
)

type companyRepository struct {
	db *gorm.DB
}

func NewCompanyRepository(db *gorm.DB) domain.CompanyRepository {
	return &companyRepository{db: db}
}

func (r *companyRepository) Create(ctx context.Context, company *domain.Company) error {
	return r.db.WithContext(ctx).Create(company).Error
}

func (r *companyRepository) CreateTx(tx interface{}, ctx context.Context, company *domain.Company) error {
	return tx.(*gorm.DB).WithContext(ctx).Create(company).Error
}

func (r *companyRepository) GetAll(ctx context.Context) ([]domain.Company, error) {
	var companies []domain.Company
	err := r.db.WithContext(ctx).Find(&companies).Error
	return companies, err
}

func (r *companyRepository) GetByEmail(ctx context.Context, email string) (*domain.Company, error) {
	var company *domain.Company
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&company).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return company, nil
}
