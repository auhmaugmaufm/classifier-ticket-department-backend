package repository

import (
	"context"

	"github.com/auhmaugmaufm/predict-ticket-department-backend/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ticketRepository struct {
	db *gorm.DB
}

func NewTicketRepositry(db *gorm.DB) domain.TicketRepository {
	return &ticketRepository{db: db}
}

func (r *ticketRepository) Create(ctx context.Context, ticket *domain.Ticket) error {
	return r.db.WithContext(ctx).Create(ticket).Error
}

func (r *ticketRepository) CreateBulk(ctx context.Context, tickets []domain.Ticket) error {
	return r.db.WithContext(ctx).CreateInBatches(tickets, 100).Error
}

func (r *ticketRepository) GetByCompanyID(ctx context.Context, company_id uuid.UUID) ([]domain.Ticket, error) {
	var ticket []domain.Ticket
	err := r.db.WithContext(ctx).Preload("Department").
		Joins("JOIN departments ON departments.id = tickets.department_id").
		Where("departments.company_id = ?", company_id).Find(&ticket).Error
	if err != nil {
		return nil, err
	}
	return ticket, nil
}
