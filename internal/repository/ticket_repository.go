package repository

import (
	"context"
	"errors"

	"github.com/auhmaugmaufm/predict-ticket-department-backend/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ticketRepository struct {
	db *gorm.DB
}

func NewTicketRepository(db *gorm.DB) domain.TicketRepository {
	return &ticketRepository{db: db}
}

func (r *ticketRepository) Create(ctx context.Context, ticket *domain.Ticket) error {
	return r.db.WithContext(ctx).Create(ticket).Error
}

func (r *ticketRepository) CreateBulk(ctx context.Context, tickets []domain.Ticket) error {
	return r.db.WithContext(ctx).CreateInBatches(tickets, 100).Error
}

func (r *ticketRepository) GetByCompanyID(ctx context.Context, companyID uuid.UUID) ([]domain.Ticket, error) {
	var tickets []domain.Ticket
	err := r.db.WithContext(ctx).Preload("Department").
		Joins("JOIN departments ON departments.id = tickets.department_id").
		Where("departments.company_id = ?", companyID).Find(&tickets).Error
	if err != nil {
		return nil, err
	}
	return tickets, nil
}

func (r *ticketRepository) UpdateTicketStatus(ctx context.Context, id uuid.UUID, status domain.TicketStatus) error {
	result := r.db.WithContext(ctx).
		Model(&domain.Ticket{}).
		Where("id = ?", id).
		Update("status", status)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("ticket not found")
	}

	return nil
}
