package handler

import (
	"context"
	"net/http"

	"github.com/auhmaugmaufm/predict-ticket-department-backend/internal/auth"
	"github.com/auhmaugmaufm/predict-ticket-department-backend/internal/domain"
	"github.com/auhmaugmaufm/predict-ticket-department-backend/internal/dto"
	"github.com/auhmaugmaufm/predict-ticket-department-backend/pkg/config"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DepartmentService interface {
	AddDepartments(ctx context.Context, departments []domain.Department) error
	GetDepartmentsByCompanyID(ctx context.Context, company_id uuid.UUID) ([]domain.Department, error)
	UpdateDepartmentStatus(ctx context.Context, id uuid.UUID, isActive bool) error
}

type DepartmentHandler struct {
	svc DepartmentService
	cfg *config.Config
}

func NewDepartmentHandler(service DepartmentService, cfg *config.Config) *DepartmentHandler {
	return &DepartmentHandler{svc: service, cfg: cfg}
}

// @Summary Create Departments
// @Description Create department in company
// @Tags department
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.DepartmentRequest true "Department credentials"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/departments [post]
func (h *DepartmentHandler) CreateDepartment(c *gin.Context) {
	var d *dto.DepartmentRequest
	if err := c.BindJSON(&d); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	companyID, ok := auth.GetCompanyID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing company_id"})
		return
	}

	departments := make([]domain.Department, len(d.DepartmentName))
	for i, departmentName := range d.DepartmentName {
		departments[i] = domain.Department{
			DepartmentName: departmentName,
			CompanyID:      companyID,
		}
	}
	err := h.svc.AddDepartments(c, departments)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "add departments success"})
}

// @Summary Get Departments By company ID
// @Description  Get Departments By company ID
// @Tags department
// @Accept json
// @Produce json
// @Security HMACAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/departments [get]
func (h *DepartmentHandler) GetDepartmentsByCompanyID(c *gin.Context) {
	companyID, err := uuid.Parse(c.Param("company_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid company_id"})
		return
	}

	departments, err := h.svc.GetDepartmentsByCompanyID(c, companyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": departments})
}

// @Summary Get Departments By company ID
// @Description  Get Departments By company ID
// @Tags department
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-HMAC-Signature header string false "HMAC signature (sha256=...)"
// @Param company_id path string true "Company ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/internal/departments/{company_id} [get]
func (h *DepartmentHandler) GetDepartmentsByCompanyIDAuth(c *gin.Context) {
	companyID, ok := auth.GetCompanyID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing company_id"})
		return
	}

	departments, err := h.svc.GetDepartmentsByCompanyID(c, companyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": departments})
}

func (h *DepartmentHandler) UpdateDepartmentStatus(c *gin.Context) {
	departmentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid ticket id",
		})
		return
	}

	var req dto.UpdateDepartmentStatusRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	err = h.svc.UpdateDepartmentStatus(c, departmentID, req.IsActive)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "department status updated",
	})
}
