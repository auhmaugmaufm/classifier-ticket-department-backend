package service

import (
	"context"
	"errors"

	"github.com/auhmaugmaufm/predict-ticket-department-backend/internal/auth"
	"github.com/auhmaugmaufm/predict-ticket-department-backend/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type CompanyService struct {
	repo           domain.CompanyRepository
	departmentRepo domain.DepartmentRepository
	txManager      domain.TxRepository
	jwtManager     *auth.JWTManager
}

func NewCompanyService(repo domain.CompanyRepository, departmentRepo domain.DepartmentRepository,
	txManager domain.TxRepository, jwtManager *auth.JWTManager) *CompanyService {
	return &CompanyService{
		repo:           repo,
		departmentRepo: departmentRepo,
		txManager:      txManager,
		jwtManager:     jwtManager,
	}
}

func (s *CompanyService) RegisterWithDefaultDepartment(ctx context.Context, email string, password string) error {
	return s.txManager.WithinTransaction(ctx, func(txCtx context.Context, tx interface{}) error {
		_, err := s.repo.GetByEmail(txCtx, email)
		if err == nil {
			return domain.ErrEmailAlreadyExists
		}
		if !errors.Is(err, domain.ErrNotFound) {
			return err
		}

		bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		company := &domain.Company{
			Email:        email,
			PasswordHash: string(bytes),
		}
		if err := s.repo.CreateTx(tx, txCtx, company); err != nil {
			return err
		}

		otherDept := &domain.Department{
			DepartmentName: "Other",
			CompanyID:      company.ID,
			IsActive:       true,
		}
		if err := s.departmentRepo.CreateTx(tx, txCtx, otherDept); err != nil {
			return err
		}

		return nil
	})
}

func (s *CompanyService) Register(ctx context.Context, email string, password string) error {
	return s.RegisterWithDefaultDepartment(ctx, email, password)
}

func (s *CompanyService) Login(ctx context.Context, email string, password string) (string, error) {
	company, err := s.repo.GetByEmail(ctx, email)
	if errors.Is(err, domain.ErrNotFound) {
		return "", domain.ErrInvalidCredentials
	}
	if err != nil {
		return "", err
	}

	comparePasswordError := bcrypt.CompareHashAndPassword([]byte(company.PasswordHash), []byte(password))
	if comparePasswordError != nil {
		return "", domain.ErrInvalidCredentials
	}

	return s.jwtManager.GenerateToken(company.ID.String(), email)
}

func (s *CompanyService) GetAllCompanies(ctx context.Context) ([]domain.Company, error) {
	return s.repo.GetAll(ctx)
}
