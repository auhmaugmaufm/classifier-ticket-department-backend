package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/auhmaugmaufm/predict-ticket-department-backend/internal/domain"
	"github.com/auhmaugmaufm/predict-ticket-department-backend/internal/dto"
	"github.com/auhmaugmaufm/predict-ticket-department-backend/pkg/config"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type FormService interface {
	SubmitForm(ctx context.Context, form *domain.Form) error
	GetSubmitFormByCompanyID(ctx context.Context, company_id uuid.UUID) ([]domain.Form, error)
	GetSubmitFormPerDayByCompanyID(ctx context.Context, company_id uuid.UUID, dateStr string) ([]domain.Form, error)
}

type FormHandler struct {
	svc FormService
	cfg *config.Config
}

func NewFormHandler(service FormService, cfg *config.Config) *FormHandler {
	return &FormHandler{svc: service, cfg: cfg}
}

// @Summary Submit Form
// @Description Submit Form to Company
// @Tags form
// @Accept json
// @Produce json
// @Param request body dto.FormRequest true "Form credentials"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/forms/submit [post]
func (h *FormHandler) SubmitForm(c *gin.Context) {
	var req *dto.FormRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	form := &domain.Form{
		LinkID:      req.LinkID,
		Title:       req.Title,
		Description: req.Description,
	}
	err := h.svc.SubmitForm(c, form)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "submit form success"})
}

// @Summary Get Forms By company ID
// @Description  Get all forms By company ID
// @Tags form
// @Accept json
// @Produce json
// @Param company_id path string true "Company ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/forms/{company_id} [get]
func (h *FormHandler) GetSubmitFormCompanyID(c *gin.Context) {
	company_id, err := uuid.Parse(c.Param("company_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid company_id"})
		return
	}
	forms, err := h.svc.GetSubmitFormByCompanyID(c, company_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": forms})
}

// @Summary Get Forms Per Day By company ID
// @Description  Get all forms per day By company ID
// @Tags form
// @Accept json
// @Produce json
// @Param company_id path string true "Company ID"
// @Param date query string false "Date in YYYY-MM-DD format; defaults to yesterday"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/forms/{company_id}/per-day [get]
func (h *FormHandler) GetSubmitFormPerDayByCompanyID(c *gin.Context) {
	company_id, err := uuid.Parse(c.Param("company_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid company_id"})
		return
	}
	dateStr := c.Query("date")
	if dateStr == "" {
		dateStr = time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	}
	forms, err := h.svc.GetSubmitFormPerDayByCompanyID(c, company_id, dateStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": forms})
}
