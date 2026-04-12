package dto

import "github.com/google/uuid"

type FormRequest struct {
	FormID      uuid.UUID `json:"form_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
}
