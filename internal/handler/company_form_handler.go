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

type CompanyFormService interface {
	CreateCompanyForm(ctx context.Context, company_id uuid.UUID) error
	GetCompanyFormByCompanyID(ctx context.Context, company_id uuid.UUID) (*domain.CompanyForm, error)
}

type CompanyFormHandler struct {
	svc CompanyFormService
	cfg *config.Config
}

func NewCompanyFormHandler(service CompanyFormService, cfg *config.Config) *CompanyFormHandler {
	return &CompanyFormHandler{svc: service, cfg: cfg}
}

// @Summary Create Company Form
// @Description Create Company Form
// @Tags company_form
// @Accept json
// @Produce json
// @Param request body dto.CompanyFormRequest true "CompanyForm credentials"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /api/v1/conmpany_form/create [post]
func (h *CompanyFormHandler) CreateCompanyForm(c *gin.Context) {
	var f *dto.CompanyFormRequest
	if err := c.BindJSON(&f); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	err := h.svc.CreateCompanyForm(c, f.CompanyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "create company form success"})
}

// @Summary Get Company Form By company ID
// @Description  Get Company Form By company ID
// @Tags company_form
// @Accept json
// @Produce json
// @Param company_id path string true "Company ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/conmpany_form/{company_id} [get]
func (h *CompanyFormHandler) GetCompanyFormByCompanyID(c *gin.Context) {
	company_id, err := uuid.Parse(c.Param("company_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid company_id"})
		return
	}

	company_form, err := h.svc.GetCompanyFormByCompanyID(c, company_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": company_form})
}
