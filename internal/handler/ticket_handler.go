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

type TicketService interface {
	CreateTicket(ctx context.Context, ticket *domain.Ticket) error
	CreateTickets(ctx context.Context, tickets []domain.Ticket) error
	GetTicketsByCompanyID(ctx context.Context, companyID uuid.UUID) ([]domain.Ticket, error)
	UpdateTicketStatusByTicketID(ctx context.Context, id uuid.UUID, status domain.TicketStatus) error
}

type TicketHandler struct {
	svc TicketService
	cfg *config.Config
}

func NewTicketHandler(service TicketService, cfg *config.Config) *TicketHandler {
	return &TicketHandler{svc: service, cfg: cfg}
}

// @Summary Create Ticket
// @Description Create a single ticket
// @Tags ticket
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.TicketRequest true "Ticket credentials"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/tickets [post]
func (h *TicketHandler) CreateTicket(c *gin.Context) {
	var req dto.TicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	priority := domain.TicketPriority(req.Priority)
	ticket := &domain.Ticket{
		Message:       req.Message,
		PredictStatus: domain.PredictStatus(req.PredictStatus),
		Status:        domain.TicketPending,
		Title:         req.Title,
		Description:   req.Description,
		DepartmentID:  req.DepartmentID,
		Priority:      &priority,
	}
	err := h.svc.CreateTicket(c, ticket)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "create ticket success"})
}

// @Summary Create Tickets
// @Description Create multiple tickets in one request
// @Tags ticket
// @Accept json
// @Produce json
// @Security BearerAuth
// @Security HMACAuth
// @Param X-HMAC-Signature header string false "HMAC signature (sha256=...)"
// @Param request body []dto.TicketRequest true "Ticket credentials"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/tickets/bulk [post]
// @Router /api/v1/internal/tickets/bulk [post]
func (h *TicketHandler) CreateTickets(c *gin.Context) {
	var req []dto.TicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	tickets := make([]domain.Ticket, len(req))
	for i, t := range req {
		priority := domain.TicketPriority(t.Priority)
		tickets[i] = domain.Ticket{
			Message:       t.Message,
			PredictStatus: domain.PredictStatus(t.PredictStatus),
			Status:        domain.TicketPending,
			Title:         t.Title,
			Description:   t.Description,
			FormID:        t.FormID,
			DepartmentID:  t.DepartmentID,
			Priority:      &priority,
		}
	}
	err := h.svc.CreateTickets(c, tickets)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "create ticket success"})
}

// @Summary Get Tickets By company ID
// @Description Get all tickets by company ID
// @Tags ticket
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/tickets [get]
func (h *TicketHandler) GetTicketsByCompanyID(c *gin.Context) {
	companyID, ok := auth.GetCompanyID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing company_id"})
		return
	}

	tickets, err := h.svc.GetTicketsByCompanyID(c, companyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": tickets})
}

// @Summary Update Ticket Status
// @Description Update status for a ticket
// @Tags ticket
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Ticket ID"
// @Param request body dto.UpdateTicketStatusRequest true "Ticket status payload"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/tickets/{id} [patch]
func (h *TicketHandler) UpdateTicketStatus(c *gin.Context) {
	ticketID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid ticket id",
		})
		return
	}

	var req dto.UpdateTicketStatusRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	err = h.svc.UpdateTicketStatusByTicketID(
		c,
		ticketID,
		req.Status,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ticket status updated",
	})
}
