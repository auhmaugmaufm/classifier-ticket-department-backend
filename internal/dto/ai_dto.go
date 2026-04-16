package dto

import (
	"github.com/auhmaugmaufm/predict-ticket-department-backend/internal/domain"
	"github.com/google/uuid"
)

type AIDataRequest struct {
	Data []CompanyFormItems `json:"data"`
}

type CompanyFormItems struct {
	CompanyID uuid.UUID     `json:"company_id"`
	Forms     []domain.Form `json:"forms"`
}

type AIResponse struct {
	Message     string `json:"message"`
	QueuedCount int    `json:"queued_count"`
}
