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
	StatusSuccess PredictStatus = "sucess"
)

const (
	PriorityLow    TicketPriority = "low"
	PriorityMedium TicketPriority = "medium"
	PriorityHigh   TicketPriority = "high"
)

type Ticket struct {
	ID            uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Message       string         `json:"message"`
	Status        PredictStatus  `json:"status"`
	Title         string         `json:"title"`
	Description   string         `json:"description"`
	DepartmentID  uuid.UUID      `json:"department_id" gorm:"type:uuid;not null;index"`
	Priority      TicketPriority `json:"priority" gorm:"not null"`
	SubmittedDate time.Time      `json:"submitted_date" gorm:"not null"`
	CreatedAt     time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	Department *Department `json:"department,omitempty" gorm:"foreignKey:DepartmentID"`
}

type TicketRepository interface {
	Create(ctx context.Context, ticket *Ticket) error
	CreateBulk(ctx context.Context, tickets []Ticket) error
	GetAll(ctx context.Context) ([]Ticket, error)
	// TODO: Get by CompanyID
}
