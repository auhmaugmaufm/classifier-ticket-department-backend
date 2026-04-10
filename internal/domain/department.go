package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Department struct {
	ID             uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	DepartmentName string         `json:"department_name" gorm:"not null"`
	CompanyID      uuid.UUID      `json:"company_id" gorm:"type:uuid;not null;index"`
	IsActive       bool           `json:"is_active" gorm:"not null;default:0"`
	CreatedAt      time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt      gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	Company *Company `json:"company,omitempty" gorm:"foreignKey:CompanyID"`
}

type DepartmentRepository interface {
	Create(ctx context.Context, department *Department) error
	GetByCompanyID(ctx context.Context, company_id uuid.UUID) ([]Department, error)
}
