package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Company struct {
	ID           uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Email        string         `json:"email" gorm:"uniqueIndex;not null"`
	PasswordHash string         `json:"-" gorm:"not null"`
	CreatedAt    time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type CompanyRepository interface {
	Create(ctx context.Context, company *Company) error
	GetAll(ctx context.Context) ([]Company, error)
	GetByEmail(ctx context.Context, email string) (*Company, error)
}
