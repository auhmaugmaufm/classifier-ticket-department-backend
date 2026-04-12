package dto

import (
	"time"

	"github.com/google/uuid"
)

type TicketRequest struct {
	Message       string    `json:"message"`
	Status        string    `json:"status"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	DepartmentID  uuid.UUID `json:"department_id"`
	Priority      string    `json:"priority"`
	SubmittedDate time.Time `json:"submitted_date"`
}
