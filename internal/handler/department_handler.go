package handler

import (
	"context"
	"net/http"

	"github.com/auhmaugmaufm/predict-ticket-department-backend/internal/domain"
	"github.com/auhmaugmaufm/predict-ticket-department-backend/internal/dto"
	"github.com/auhmaugmaufm/predict-ticket-department-backend/pkg/config"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DepartmentService interface {
	AddDepartments(ctx context.Context, departments []domain.Department) error
	GetDepartmentsByCompanyID(ctx context.Context, company_id uuid.UUID) ([]domain.Department, error)
}

type DeparmtentHandler struct {
	svc DepartmentService
	cfg *config.Config
}

func NewDepartmentHandler(service DepartmentService, cfg *config.Config) *DeparmtentHandler {
	return &DeparmtentHandler{svc: service, cfg: cfg}
}

// @Summary Add Departments
// @Description Create department in company
// @Tags department
// @Accept json
// @Produce json
// @Param request body dto.DepartmentRequest true "Department credentials"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /api/v1/departments/add [post]
func (h *DeparmtentHandler) AddDepartments(c *gin.Context) {
	var d []dto.DepartmentRequest
	if err := c.BindJSON(&d); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	departments := make([]domain.Department, len(d))
	for i, department := range d {
		departments[i] = domain.Department{
			DepartmentName: department.DepartmentName,
			CompanyID:      department.CompanyID,
		}
	}
	err := h.svc.AddDepartments(c, departments)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "add departments success"})
}

// @Summary Get Departments By company ID
// @Description  Get Departments By company ID
// @Tags department
// @Accept json
// @Produce json
// @Param company_id path string true "Company ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/departments/company/{company_id} [get]
func (h *DeparmtentHandler) GetDepartmentsByCompanyID(c *gin.Context) {
	company_id, err := uuid.Parse(c.Param("company_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid company_id"})
		return
	}

	departments, err := h.svc.GetDepartmentsByCompanyID(c, company_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": departments})
}
