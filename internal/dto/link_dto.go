package dto

import "github.com/google/uuid"

type LinkRequest struct {
	CompanyID uuid.UUID `json:"company_id"`
}
