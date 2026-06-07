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

type MockFormRepo struct{ mock.Mock }

func (m *MockFormRepo) Create(ctx context.Context, form *domain.Form) error {
	args := m.Called(ctx, form)
	return args.Error(0)
}

func (m *MockFormRepo) GetByCompanyID(ctx context.Context, company_id uuid.UUID) ([]domain.Form, error) {
	args := m.Called(ctx, company_id)
	return args.Get(0).([]domain.Form), args.Error(1)
}

func (m *MockFormRepo) GetFormCompanyID(ctx context.Context, company_id uuid.UUID, dateStr string) ([]domain.Form, error) {
	args := m.Called(ctx, company_id, dateStr)
	return args.Get(0).([]domain.Form), args.Error(1)
}

func TestSubmitForm_Success(t *testing.T) {
	repo := new(MockFormRepo)
	form := &domain.Form{CompanyID: uuid.New(), Title: "Bug Report", Description: "desc"}
	repo.On("Create", mock.Anything, form).Return(nil)

	svc := NewFormService(repo)
	err := svc.SubmitForm(context.Background(), form)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestSubmitForm_Error(t *testing.T) {
	repo := new(MockFormRepo)
	form := &domain.Form{CompanyID: uuid.New(), Title: "Bug Report", Description: "desc"}
	repo.On("Create", mock.Anything, form).Return(errors.New("db error"))

	svc := NewFormService(repo)
	err := svc.SubmitForm(context.Background(), form)

	assert.Error(t, err)
	repo.AssertExpectations(t)
}

func TestGetSubmitFormByCompanyID_Success(t *testing.T) {
	repo := new(MockFormRepo)
	companyID := uuid.New()
	expected := []domain.Form{{Title: "Form A"}, {Title: "Form B"}}
	repo.On("GetByCompanyID", mock.Anything, companyID).Return(expected, nil)

	svc := NewFormService(repo)
	result, err := svc.GetSubmitFormByCompanyID(context.Background(), companyID)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	repo.AssertExpectations(t)
}

func TestGetSubmitFormByCompanyID_Error(t *testing.T) {
	repo := new(MockFormRepo)
	companyID := uuid.New()
	repo.On("GetByCompanyID", mock.Anything, companyID).
		Return([]domain.Form{}, errors.New("db error"))

	svc := NewFormService(repo)
	result, err := svc.GetSubmitFormByCompanyID(context.Background(), companyID)

	assert.Error(t, err)
	assert.Empty(t, result)
	repo.AssertExpectations(t)
}

func TestGetSubmitFormPerDayByCompanyID_Success(t *testing.T) {
	repo := new(MockFormRepo)
	companyID := uuid.New()
	dateStr := "2025-01-15"
	expected := []domain.Form{{Title: "Daily Form"}}
	repo.On("GetFormCompanyID", mock.Anything, companyID, dateStr).Return(expected, nil)

	svc := NewFormService(repo)
	result, err := svc.GetSubmitFormPerDayByCompanyID(context.Background(), companyID, dateStr)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	repo.AssertExpectations(t)
}

func TestGetSubmitFormPerDayByCompanyID_Error(t *testing.T) {
	repo := new(MockFormRepo)
	companyID := uuid.New()
	dateStr := "2025-01-15"
	repo.On("GetFormCompanyID", mock.Anything, companyID, dateStr).
		Return([]domain.Form{}, errors.New("db error"))

	svc := NewFormService(repo)
	result, err := svc.GetSubmitFormPerDayByCompanyID(context.Background(), companyID, dateStr)

	assert.Error(t, err)
	assert.Empty(t, result)
	repo.AssertExpectations(t)
}
