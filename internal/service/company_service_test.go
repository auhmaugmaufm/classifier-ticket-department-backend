package service

import (
	"context"
	"errors"
	"testing"

	"github.com/auhmaugmaufm/predict-ticket-department-backend/internal/auth"
	"github.com/auhmaugmaufm/predict-ticket-department-backend/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

type MockCompanyRepo struct {
	mock.Mock
}

func (m *MockCompanyRepo) GetByEmail(ctx context.Context, email string) (*domain.Company, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Company), args.Error(1)
}

func (m *MockCompanyRepo) CreateTx(tx interface{}, ctx context.Context, c *domain.Company) error {
	args := m.Called(tx, ctx, c)
	return args.Error(0)
}

func (m *MockCompanyRepo) GetAll(ctx context.Context) ([]domain.Company, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.Company), args.Error(1)
}

type MockTxManager struct{ mock.Mock }

func (m *MockTxManager) WithinTransaction(ctx context.Context, fn func(context.Context, interface{}) error) error {
	return fn(ctx, nil)
}

func TestRegister_Success(t *testing.T) {
	companyRepo := new(MockCompanyRepo)
	deptRepo := new(MockDepartmentRepo)
	txMgr := new(MockTxManager)
	jwtMgr := auth.NewJWTManager("secret", 24)

	companyRepo.On("GetByEmail", mock.Anything, "new@test.com").
		Return(nil, domain.ErrNotFound)
	companyRepo.On("CreateTx", mock.Anything, mock.Anything, mock.Anything).
		Return(nil)
	deptRepo.On("CreateTx", mock.Anything, mock.Anything, mock.Anything).
		Return(nil)

	svc := NewCompanyService(companyRepo, deptRepo, txMgr, jwtMgr)
	err := svc.Register(context.Background(), "new@test.com", "password123")

	assert.NoError(t, err)
	companyRepo.AssertExpectations(t)
	deptRepo.AssertExpectations(t)
}

func TestRegister_EmailAlreadyExists(t *testing.T) {
	companyRepo := new(MockCompanyRepo)
	deptRepo := new(MockDepartmentRepo)
	txMgr := new(MockTxManager)
	jwtMgr := auth.NewJWTManager("secret", 24)

	companyRepo.On("GetByEmail", mock.Anything, "exist@test.com").
		Return(&domain.Company{Email: "exist@test.com"}, nil)

	svc := NewCompanyService(companyRepo, deptRepo, txMgr, jwtMgr)
	err := svc.Register(context.Background(), "exist@test.com", "password123")

	assert.ErrorIs(t, err, domain.ErrEmailAlreadyExists)
	companyRepo.AssertNotCalled(t, "CreateTx")
}

func TestRegister_RepoError(t *testing.T) {
	companyRepo := new(MockCompanyRepo)
	deptRepo := new(MockDepartmentRepo)
	txMgr := new(MockTxManager)
	jwtMgr := auth.NewJWTManager("secret", 24)

	companyRepo.On("GetByEmail", mock.Anything, "x@test.com").
		Return(nil, errors.New("db error"))

	svc := NewCompanyService(companyRepo, deptRepo, txMgr, jwtMgr)
	err := svc.Register(context.Background(), "x@test.com", "pass")

	assert.Error(t, err)
	assert.NotErrorIs(t, err, domain.ErrEmailAlreadyExists)
}

func TestLogin_Success(t *testing.T) {
	companyRepo := new(MockCompanyRepo)
	deptRepo := new(MockDepartmentRepo)
	txMgr := new(MockTxManager)
	jwtMgr := auth.NewJWTManager("secret", 24)

	hashed, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	company := &domain.Company{
		ID:           uuid.New(),
		Email:        "user@test.com",
		PasswordHash: string(hashed),
	}
	companyRepo.On("GetByEmail", mock.Anything, "user@test.com").Return(company, nil)

	svc := NewCompanyService(companyRepo, deptRepo, txMgr, jwtMgr)
	token, err := svc.Login(context.Background(), "user@test.com", "password123")

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestLogin_UserNotFound(t *testing.T) {
	companyRepo := new(MockCompanyRepo)
	deptRepo := new(MockDepartmentRepo)
	txMgr := new(MockTxManager)
	jwtMgr := auth.NewJWTManager("secret", 24)

	companyRepo.On("GetByEmail", mock.Anything, "ghost@test.com").
		Return(nil, domain.ErrNotFound)

	svc := NewCompanyService(companyRepo, deptRepo, txMgr, jwtMgr)
	token, err := svc.Login(context.Background(), "ghost@test.com", "pass")

	assert.Empty(t, token)
	assert.ErrorIs(t, err, domain.ErrInvalidCredentials)
}

func TestLogin_WrongPassword(t *testing.T) {
	companyRepo := new(MockCompanyRepo)
	deptRepo := new(MockDepartmentRepo)
	txMgr := new(MockTxManager)
	jwtMgr := auth.NewJWTManager("secret", 24)

	hashed, _ := bcrypt.GenerateFromPassword([]byte("correctpass"), bcrypt.DefaultCost)
	company := &domain.Company{
		ID:           uuid.New(),
		Email:        "user@test.com",
		PasswordHash: string(hashed),
	}
	companyRepo.On("GetByEmail", mock.Anything, "user@test.com").Return(company, nil)

	svc := NewCompanyService(companyRepo, deptRepo, txMgr, jwtMgr)
	token, err := svc.Login(context.Background(), "user@test.com", "wrongpass")

	assert.Empty(t, token)
	assert.ErrorIs(t, err, domain.ErrInvalidCredentials)
}

func TestGetAllCompanies_Success(t *testing.T) {
	companyRepo := new(MockCompanyRepo)
	deptRepo := new(MockDepartmentRepo)
	txMgr := new(MockTxManager)
	jwtMgr := auth.NewJWTManager("secret", 24)

	companyRepo.On("GetAll", mock.Anything).
		Return([]domain.Company{{Email: "a@test.com"}, {Email: "b@test.com"}}, nil)

	svc := NewCompanyService(companyRepo, deptRepo, txMgr, jwtMgr)
	companies, err := svc.GetAllCompanies(context.Background())

	assert.NoError(t, err)
	assert.Len(t, companies, 2)
}

func TestGetAllCompanies_Error(t *testing.T) {
	companyRepo := new(MockCompanyRepo)
	deptRepo := new(MockDepartmentRepo)
	txMgr := new(MockTxManager)
	jwtMgr := auth.NewJWTManager("secret", 24)

	companyRepo.On("GetAll", mock.Anything).
		Return([]domain.Company{}, errors.New("db error"))

	svc := NewCompanyService(companyRepo, deptRepo, txMgr, jwtMgr)
	companies, err := svc.GetAllCompanies(context.Background())

	assert.Error(t, err)
	assert.Empty(t, companies)
}
