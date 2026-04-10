package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Form struct {
	ID            uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	FormID        uuid.UUID      `json:"form_id" gorm:"type:uuid;not null;index"`
	Title         string         `json:"title" gorm:"not null"`
	Description   string         `json:"description" gorm:"not null"`
	SubmittedDate time.Time      `json:"submitted_date" gorm:"not null"`
	CreatedAt     time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type FormRepository interface {
	Create(ctx context.Context, form *Form) error
	GetByFormID(ctx context.Context, form_id uuid.UUID) ([]Form, error)
}
