package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Link struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CompanyID uuid.UUID      `json:"company_id" gorm:"type:uuid;not null;index"`
	Link      string         `json:"link_form" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	Company *Company `json:"company,omitempty" gorm:"foreignKey:CompanyID"`
}

type LinkRepository interface {
	Create(ctx context.Context, Link *Link) error
	GetByID(ctx context.Context, company_id uuid.UUID) (*Link, error)
}
