package service

import (
	"context"
	"errors"

	"github.com/auhmaugmaufm/predict-ticket-department-backend/internal/auth"
	"github.com/auhmaugmaufm/predict-ticket-department-backend/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type CompanyService struct {
	repo       domain.CompanyRepository
	jwtManager *auth.JWTManager
}

func NewCompanyService(repo domain.CompanyRepository, jwtManager *auth.JWTManager) *CompanyService {
	return &CompanyService{repo: repo, jwtManager: jwtManager}
}

func (s *CompanyService) Register(ctx context.Context, email string, password string) error {
	_, err := s.repo.GetByEmail(ctx, email)
	if err == nil {
		return errors.New("Email already exists")
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	company := &domain.Company{
		Email:        email,
		PasswordHash: string(bytes),
	}
	return s.repo.Create(ctx, company)
}

func (s *CompanyService) Login(ctx context.Context, email string, password string) (string, error) {
	company, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	comparePasswordError := bcrypt.CompareHashAndPassword([]byte(company.PasswordHash), []byte(password))
	if comparePasswordError != nil {
		return "", comparePasswordError
	}

	return s.jwtManager.GenerateToken(company.ID.String(), email)
}
