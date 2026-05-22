package dto

import "github.com/google/uuid"

type FormRequest struct {
	CompanyID   uuid.UUID `json:"company_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
}
