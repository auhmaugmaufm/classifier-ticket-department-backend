package dto

import (
	"github.com/google/uuid"
)

type TicketRequest struct {
	Message      string    `json:"message"`
	Status       string    `json:"status"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	FormID       uuid.UUID `json:"form_id"`
	DepartmentID uuid.UUID `json:"department_id"`
	Priority     string    `json:"priority"`
}
