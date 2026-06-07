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

type MockTicketRepo struct{ mock.Mock }

func (m *MockTicketRepo) Create(ctx context.Context, ticket *domain.Ticket) error {
	args := m.Called(ctx, ticket)
	return args.Error(0)
}

func (m *MockTicketRepo) CreateBulk(ctx context.Context, tickets []domain.Ticket) error {
	args := m.Called(ctx, tickets)
	return args.Error(0)
}

func (m *MockTicketRepo) GetByCompanyID(ctx context.Context, companyID uuid.UUID) ([]domain.Ticket, error) {
	args := m.Called(ctx, companyID)
	return args.Get(0).([]domain.Ticket), args.Error(1)
}

func (m *MockTicketRepo) UpdateTicketStatus(ctx context.Context, id uuid.UUID, status domain.TicketStatus) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func TestCreateTicket_Success(t *testing.T) {
	repo := new(MockTicketRepo)
	ticket := &domain.Ticket{FormID: uuid.New(), Title: "Bug", Message: "msg"}
	repo.On("Create", mock.Anything, ticket).Return(nil)

	svc := NewTicketService(repo)
	err := svc.CreateTicket(context.Background(), ticket)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestCreateTicket_Error(t *testing.T) {
	repo := new(MockTicketRepo)
	ticket := &domain.Ticket{FormID: uuid.New(), Title: "Bug", Message: "msg"}
	repo.On("Create", mock.Anything, ticket).Return(errors.New("db error"))

	svc := NewTicketService(repo)
	err := svc.CreateTicket(context.Background(), ticket)

	assert.Error(t, err)
	repo.AssertExpectations(t)
}

func TestCreateTickets_Success(t *testing.T) {
	repo := new(MockTicketRepo)
	tickets := []domain.Ticket{
		{FormID: uuid.New(), Title: "T1"},
		{FormID: uuid.New(), Title: "T2"},
	}
	repo.On("CreateBulk", mock.Anything, tickets).Return(nil)

	svc := NewTicketService(repo)
	err := svc.CreateTickets(context.Background(), tickets)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestCreateTickets_Error(t *testing.T) {
	repo := new(MockTicketRepo)
	tickets := []domain.Ticket{{FormID: uuid.New(), Title: "T1"}}
	repo.On("CreateBulk", mock.Anything, tickets).Return(errors.New("db error"))

	svc := NewTicketService(repo)
	err := svc.CreateTickets(context.Background(), tickets)

	assert.Error(t, err)
	repo.AssertExpectations(t)
}

func TestGetTicketsByCompanyID_Success(t *testing.T) {
	repo := new(MockTicketRepo)
	companyID := uuid.New()
	expected := []domain.Ticket{{Title: "T1"}, {Title: "T2"}}
	repo.On("GetByCompanyID", mock.Anything, companyID).Return(expected, nil)

	svc := NewTicketService(repo)
	result, err := svc.GetTicketsByCompanyID(context.Background(), companyID)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	repo.AssertExpectations(t)
}

func TestGetTicketsByCompanyID_Error(t *testing.T) {
	repo := new(MockTicketRepo)
	companyID := uuid.New()
	repo.On("GetByCompanyID", mock.Anything, companyID).
		Return([]domain.Ticket{}, errors.New("db error"))

	svc := NewTicketService(repo)
	result, err := svc.GetTicketsByCompanyID(context.Background(), companyID)

	assert.Error(t, err)
	assert.Empty(t, result)
	repo.AssertExpectations(t)
}

func TestUpdateTicketStatusByTicketID_Success(t *testing.T) {
	repo := new(MockTicketRepo)
	id := uuid.New()
	repo.On("UpdateTicketStatus", mock.Anything, id, domain.TicketClosed).Return(nil)

	svc := NewTicketService(repo)
	err := svc.UpdateTicketStatusByTicketID(context.Background(), id, domain.TicketClosed)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestUpdateTicketStatusByTicketID_Error(t *testing.T) {
	repo := new(MockTicketRepo)
	id := uuid.New()
	repo.On("UpdateTicketStatus", mock.Anything, id, domain.TicketPending).
		Return(errors.New("db error"))

	svc := NewTicketService(repo)
	err := svc.UpdateTicketStatusByTicketID(context.Background(), id, domain.TicketPending)

	assert.Error(t, err)
	repo.AssertExpectations(t)
}
