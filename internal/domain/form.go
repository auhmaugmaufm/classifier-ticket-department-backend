package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Form struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	FormID      uuid.UUID      `json:"form_id" gorm:"type:uuid;not null;index"`
	Title       string         `json:"title" gorm:"not null"`
	Description string         `json:"description" gorm:"not null"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type FormRepository interface {
	Create(ctx context.Context, form *Form) error
	GetByCompanyID(ctx context.Context, company_id uuid.UUID) ([]Form, error)
	GetPerDayByCompanyID(ctx context.Context, company_id uuid.UUID, dateStr string) ([]Form, error)
}
