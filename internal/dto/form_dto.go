package dto

import "github.com/google/uuid"

type FormRequest struct {
	LinkID      uuid.UUID `json:"link_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
}
