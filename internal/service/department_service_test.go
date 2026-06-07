package service

import (
	"context"
	"errors"
	"testing"

	"github.com/auhmaugmaufm/predict-ticket-department-backend/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDepartmentRepo struct{ mock.Mock }

func (m *MockDepartmentRepo) Create(ctx context.Context, department *domain.Department) error {
	args := m.Called(ctx, department)
	return args.Error(0)
}

func (m *MockDepartmentRepo) CreateBulk(ctx context.Context, departments []domain.Department) error {
	args := m.Called(ctx, departments)
	return args.Error(0)
}

func (m *MockDepartmentRepo) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockDepartmentRepo) GetByCompanyID(ctx context.Context, company_id uuid.UUID) ([]domain.Department, error) {
	args := m.Called(ctx, company_id)
	return args.Get(0).([]domain.Department), args.Error(1)
}

func (m *MockDepartmentRepo) UpdateStatus(ctx context.Context, id uuid.UUID, isActive bool) error {
	args := m.Called(ctx, id, isActive)
	return args.Error(0)
}

func (m *MockDepartmentRepo) CreateTx(tx interface{}, ctx context.Context, d *domain.Department) error {
	args := m.Called(tx, ctx, d)
	return args.Error(0)
}

func TestAddDepartments_Success(t *testing.T) {
	repo := new(MockDepartmentRepo)
	departments := []domain.Department{
		{DepartmentName: "IT", CompanyID: uuid.New()},
		{DepartmentName: "HR", CompanyID: uuid.New()},
	}
	repo.On("CreateBulk", mock.Anything, departments).Return(nil)

	svc := NewDepartmentService(repo)
	err := svc.AddDepartments(context.Background(), departments)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestAddDepartments_Error(t *testing.T) {
	repo := new(MockDepartmentRepo)
	repo.On("CreateBulk", mock.Anything, mock.Anything).Return(errors.New("db error"))

	svc := NewDepartmentService(repo)
	err := svc.AddDepartments(context.Background(), []domain.Department{})

	assert.Error(t, err)
	repo.AssertExpectations(t)
}

func TestGetDepartmentsByCompanyID_Success(t *testing.T) {
	repo := new(MockDepartmentRepo)
	companyID := uuid.New()
	expected := []domain.Department{
		{DepartmentName: "IT", CompanyID: companyID},
	}
	repo.On("GetByCompanyID", mock.Anything, companyID).Return(expected, nil)

	svc := NewDepartmentService(repo)
	result, err := svc.GetDepartmentsByCompanyID(context.Background(), companyID)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	repo.AssertExpectations(t)
}

func TestGetDepartmentsByCompanyID_Error(t *testing.T) {
	repo := new(MockDepartmentRepo)
	companyID := uuid.New()
	repo.On("GetByCompanyID", mock.Anything, companyID).
		Return([]domain.Department{}, errors.New("db error"))

	svc := NewDepartmentService(repo)
	result, err := svc.GetDepartmentsByCompanyID(context.Background(), companyID)

	assert.Error(t, err)
	assert.Empty(t, result)
	repo.AssertExpectations(t)
}

func TestUpdateDepartmentStatus_Success(t *testing.T) {
	repo := new(MockDepartmentRepo)
	id := uuid.New()
	repo.On("UpdateStatus", mock.Anything, id, true).Return(nil)

	svc := NewDepartmentService(repo)
	err := svc.UpdateDepartmentStatus(context.Background(), id, true)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestUpdateDepartmentStatus_Error(t *testing.T) {
	repo := new(MockDepartmentRepo)
	id := uuid.New()
	repo.On("UpdateStatus", mock.Anything, id, false).Return(errors.New("db error"))

	svc := NewDepartmentService(repo)
	err := svc.UpdateDepartmentStatus(context.Background(), id, false)

	assert.Error(t, err)
	repo.AssertExpectations(t)
}

func TestDeleteDepartmentByID_Success(t *testing.T) {
	repo := new(MockDepartmentRepo)
	id := uuid.New()
	repo.On("Delete", mock.Anything, id).Return(nil)

	svc := NewDepartmentService(repo)
	err := svc.DeleteDepartmentByID(context.Background(), id)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestDeleteDepartmentByID_Error(t *testing.T) {
	repo := new(MockDepartmentRepo)
	id := uuid.New()
	repo.On("Delete", mock.Anything, id).Return(errors.New("db error"))

	svc := NewDepartmentService(repo)
	err := svc.DeleteDepartmentByID(context.Background(), id)

	assert.Error(t, err)
	repo.AssertExpectations(t)
}
