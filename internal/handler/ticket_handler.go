package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/auhmaugmaufm/predict-ticket-department-backend/internal/domain"
	"github.com/auhmaugmaufm/predict-ticket-department-backend/internal/dto"
	"github.com/auhmaugmaufm/predict-ticket-department-backend/pkg/config"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TicketService interface {
	CreateTicket(ctx context.Context, ticket *domain.Ticket) error
	CreateTickets(ctx context.Context, tickets []domain.Ticket) error
	GetTicketsByCompanyID(ctx context.Context, company_id uuid.UUID) ([]domain.Ticket, error)
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
// @Param request body dto.TicketRequest true "Ticket credentials"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/tickets/create [post]
func (h *TicketHandler) CreateTicket(c *gin.Context) {
	var req *dto.TicketRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	priority := domain.TicketPriority(req.Priority)
	ticket := &domain.Ticket{
		Message:      req.Message,
		Status:       domain.PredictStatus(req.Status),
		Title:        req.Title,
		Description:  req.Description,
		DepartmentID: req.DepartmentID,
		Priority:     &priority,
	}
	err := h.svc.CreateTicket(c, ticket)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "create ticket success"})
}

// @Summary Create Tickets
// @Description Create multiple tickets in one request
// @Tags ticket
// @Accept json
// @Produce json
// @Param request body []dto.TicketRequest true "Ticket credentials"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/tickets/create-bulk [post]
func (h *TicketHandler) CreateTickets(c *gin.Context) {
	var req []dto.TicketRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	fmt.Printf("%v\n", req)
	tickets := make([]domain.Ticket, len(req))
	for i, t := range req {
		priority := domain.TicketPriority(t.Priority)
		tickets[i] = domain.Ticket{
			Message:      t.Message,
			Status:       domain.PredictStatus(t.Status),
			Title:        t.Title,
			Description:  t.Description,
			FormID:       t.FormID,
			DepartmentID: t.DepartmentID,
			Priority:     &priority,
		}
	}
	err := h.svc.CreateTickets(c, tickets)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "create ticket success"})
}

// @Summary Get Tickets By company ID
// @Description Get all tickets by company ID
// @Tags ticket
// @Accept json
// @Produce json
// @Param company_id path string true "Company ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/tickets/{company_id} [get]
func (h *TicketHandler) GetTicketsByCompanyID(c *gin.Context) {
	company_id, err := uuid.Parse(c.Param("company_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid company_id"})
		return
	}

	tickets, err := h.svc.GetTicketsByCompanyID(c, company_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": tickets})
}
