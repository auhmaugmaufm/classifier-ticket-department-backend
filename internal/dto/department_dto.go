package dto

import (
	"github.com/google/uuid"
)

type UpdateDepartmentStatusRequest struct {
	IsActive bool `json:"is_active"`
}

type DepartmentRequest struct {
	DepartmentName []string `json:"department_name"`
}

type DepartmentResponse struct {
	ID             uuid.UUID `json:"department_id"`
	DepartmentName string    `json:"department_name"`
}
