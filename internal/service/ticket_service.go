package service

import (
	"context"

	"github.com/auhmaugmaufm/predict-ticket-department-backend/internal/domain"
	"github.com/google/uuid"
)

type TicketService struct {
	repo domain.TicketRepository
}

func NewTicketService(repo domain.TicketRepository) *TicketService {
	return &TicketService{repo: repo}
}

func (s *TicketService) CreateTicket(ctx context.Context, ticket *domain.Ticket) error {
	return s.repo.Create(ctx, ticket)
}

func (s *TicketService) CreateTickets(ctx context.Context, tickets []domain.Ticket) error {
	return s.repo.CreateBulk(ctx, tickets)
}

func (s *TicketService) GetTicketsByCompanyID(ctx context.Context, companyID uuid.UUID) ([]domain.Ticket, error) {
	return s.repo.GetByCompanyID(ctx, companyID)
}

func (s *TicketService) UpdateTicketStatusByTicketID(ctx context.Context, id uuid.UUID, status domain.TicketStatus) error {
	return s.repo.UpdateTicketStatus(ctx, id, status)
}
