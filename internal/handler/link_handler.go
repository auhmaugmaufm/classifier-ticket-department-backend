package handler

import (
	"context"
	"net/http"

	"github.com/auhmaugmaufm/predict-ticket-department-backend/internal/auth"
	"github.com/auhmaugmaufm/predict-ticket-department-backend/internal/domain"
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
	companyID, ok := auth.GetCompanyID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing company_id"})
		return
	}
	err := h.svc.CreateLink(c, companyID)
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
	companyID, ok := auth.GetCompanyID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing company_id"})
		return
	}

	company_form, err := h.svc.GetLinkByCompanyID(c, companyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": company_form})
}
