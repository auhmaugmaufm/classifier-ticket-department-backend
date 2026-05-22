package dto

import (
	"github.com/auhmaugmaufm/predict-ticket-department-backend/internal/domain"
	"github.com/google/uuid"
)

type UpdateTicketStatusRequest struct {
	Status domain.TicketStatus `json:"status"`
}

type TicketRequest struct {
	Message       string     `json:"message"`
	PredictStatus string     `json:"predict_status"`
	Title         string     `json:"title"`
	Description   string     `json:"description"`
	FormID        uuid.UUID  `json:"form_id"`
	DepartmentID  *uuid.UUID `json:"department_id"`
	Priority      string     `json:"priority"`
}
