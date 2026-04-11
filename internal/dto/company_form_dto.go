package dto

import "github.com/google/uuid"

type CompanyFormRequest struct {
	CompanyID uuid.UUID `json:"company_id"`
}
