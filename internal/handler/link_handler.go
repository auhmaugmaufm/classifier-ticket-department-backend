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
// @Tags link
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/links [post]
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
// @Tags link
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/links [get]
func (h *LinkHandler) GetLinkByCompanyID(c *gin.Context) {
	companyID, ok := auth.GetCompanyID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing company_id"})
		return
	}

	companyForm, err := h.svc.GetLinkByCompanyID(c, companyID)
	if companyForm == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "link not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": companyForm})
}
