package dto

import (
	"github.com/google/uuid"
)

type UpdateDepartmentStatusRequest struct {
	IsActive *bool `json:"is_active" binding:"required"`
}

type DepartmentRequest struct {
	DepartmentName []string `json:"department_name" binding:"required,min=1,dive,required"`
}

type DepartmentResponse struct {
	ID             uuid.UUID `json:"department_id"`
	DepartmentName string    `json:"department_name"`
}
