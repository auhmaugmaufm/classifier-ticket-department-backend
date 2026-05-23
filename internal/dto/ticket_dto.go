package dto

import (
	"github.com/auhmaugmaufm/predict-ticket-department-backend/internal/domain"
	"github.com/google/uuid"
)

type UpdateTicketStatusRequest struct {
	Status domain.TicketStatus `json:"status" binding:"required,oneof=opened pending closed"`
}

type TicketRequest struct {
	Message       string     `json:"message"        binding:"required"`
	PredictStatus string     `json:"predict_status" binding:"required,oneof=failed success"`
	Title         string     `json:"title"          binding:"required"`
	Description   string     `json:"description"    binding:"required"`
	FormID        uuid.UUID  `json:"form_id"        binding:"required"`
	DepartmentID  *uuid.UUID `json:"department_id"`
	Priority      string     `json:"priority"       binding:"required,oneof=low medium high"`
}
