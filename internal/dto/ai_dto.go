package dto

import (
	"github.com/google/uuid"
)

type AIDataRequest struct {
	Data []CompanyFormItems `json:"data"`
}

type AIForm struct {
	ID uuid.UUID `json:"id"`
	// LinkID      uuid.UUID `json:"link_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type CompanyFormItems struct {
	CompanyID uuid.UUID `json:"company_id"`
	Forms     []AIForm  `json:"forms"`
}

type AIResponse struct {
	Message     string `json:"message"`
	QueuedCount int    `json:"queued_count"`
}
