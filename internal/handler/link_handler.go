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

type LinkService interface {
	CreateLink(ctx context.Context, company_id uuid.UUID) error
	GetLinkByCompanyID(ctx context.Context, company_id uuid.UUID) (*domain.Link, error)
}

type LinkHandler struct {
	svc LinkService
	cfg *config.Config
}

func NewLinkHandler(service LinkService, cfg *config.Config) *LinkHandler {
	return &LinkHandler{svc: service, cfg: cfg}
}

// @Summary Create Company Form
// @Description Create Company Form
// @Tags company_form
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.LinkRequest true "Link credentials"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/company_form/create [post]
func (h *LinkHandler) CreateLink(c *gin.Context) {
	var f *dto.LinkRequest
	if err := c.BindJSON(&f); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	err := h.svc.CreateLink(c, f.CompanyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "created company form success"})
}

// @Summary Get Company Form By company ID
// @Description  Get Company Form By company ID
// @Tags company_form
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param company_id path string true "Company ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/company_form/{company_id} [get]
func (h *LinkHandler) GetLinkByCompanyID(c *gin.Context) {
	company_id, err := uuid.Parse(c.Param("company_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid company_id"})
		return
	}

	company_form, err := h.svc.GetLinkByCompanyID(c, company_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": company_form})
}
