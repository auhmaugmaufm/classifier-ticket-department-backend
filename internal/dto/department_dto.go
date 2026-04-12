package dto

import "github.com/google/uuid"

type DepartmentRequest struct {
	DepartmentName []string  `json:"department_name"`
	CompanyID      uuid.UUID `json:"company_id"`
}

type DepartmentResponse struct {
	ID             uuid.UUID `json:"department_id"`
	DepartmentName string    `json:"department_name"`
	CompanyID      uuid.UUID `json:"company_id"`
}
