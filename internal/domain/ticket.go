package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PredictStatus string
type TicketPriority string

const (
	StatusFailed  PredictStatus = "failed"
	StatusSuccess PredictStatus = "success"
)

const (
	PriorityLow    TicketPriority = "low"
	PriorityMedium TicketPriority = "medium"
	PriorityHigh   TicketPriority = "high"
)

type Ticket struct {
	ID           uuid.UUID       `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Message      string          `json:"message"`
	Status       PredictStatus   `json:"status"`
	Title        string          `json:"title"`
	Description  string          `json:"description"`
	FormID       uuid.UUID       `json:"form_id" gorm:"type:uuid;not null;index"`
	DepartmentID *uuid.UUID      `json:"department_id,omitempty" gorm:"type:uuid;index"`
	Priority     *TicketPriority `json:"priority,omitempty"`
	CreatedAt    time.Time       `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time       `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt  `json:"deleted_at" gorm:"index"`

	Form       *Form       `json:"form,omitempty" gorm:"foreignKey:FormID"`
	Department *Department `json:"department,omitempty" gorm:"foreignKey:DepartmentID"`
}

type TicketRepository interface {
	Create(ctx context.Context, ticket *Ticket) error
	CreateBulk(ctx context.Context, tickets []Ticket) error
	GetByCompanyID(ctx context.Context, company_id uuid.UUID) ([]Ticket, error)
}
