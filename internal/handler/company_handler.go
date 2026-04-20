package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/auhmaugmaufm/predict-ticket-department-backend/internal/domain"
	"github.com/auhmaugmaufm/predict-ticket-department-backend/internal/dto"
	"github.com/auhmaugmaufm/predict-ticket-department-backend/pkg/config"
	"github.com/gin-gonic/gin"
)

type CompanyService interface {
	Login(ctx context.Context, email string, password string) (string, error)
	Register(ctx context.Context, email string, password string) error
}

type CompanyHandler struct {
	svc CompanyService
	cfg *config.Config
}

func NewCompanyHandler(service CompanyService, cfg *config.Config) *CompanyHandler {
	return &CompanyHandler{svc: service, cfg: cfg}
}

// @Summary Register company
// @Description Create new company accout
// @Tags Company
// @Accept json
// @Produce json
// @Param request body dto.CompanyRequest true "Company credentials"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/register [post]
func (h *CompanyHandler) Register(c *gin.Context) {
	var company *dto.CompanyRequest
	if err := c.BindJSON(&company); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	err := h.svc.Register(c, company.Email, company.Password)
	if err != nil {
		if errors.Is(err, domain.ErrEmailAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "register company success"})
}

// @Summary Login company
// @Description Authenticate and get JWT token
// @Tags Company
// @Accept json
// @Produce json
// @Param request body dto.CompanyRequest true "Company credentials"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/v1/login [post]
func (h *CompanyHandler) Login(c *gin.Context) {
	var company *dto.CompanyRequest
	if err := c.BindJSON(&company); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	token, err := h.svc.Login(c, company.Email, company.Password)
	if errors.Is(err, domain.ErrInvalidCredentials) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}
