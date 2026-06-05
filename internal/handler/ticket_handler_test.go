package handler

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/auhmaugmaufm/predict-ticket-department-backend/internal/domain"
	"github.com/auhmaugmaufm/predict-ticket-department-backend/pkg/config"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTicketService struct {
	mock.Mock
}

func (m *MockTicketService) CreateTicket(ctx context.Context, ticket *domain.Ticket) error {
	args := m.Called(ctx, ticket)
	return args.Error(0)
}

func (m *MockTicketService) CreateTickets(ctx context.Context, tickets []domain.Ticket) error {
	args := m.Called(ctx, tickets)
	return args.Error(0)
}

func (m *MockTicketService) GetTicketsByCompanyID(ctx context.Context, companyID uuid.UUID) ([]domain.Ticket, error) {
	args := m.Called(ctx, companyID)
	return args.Get(0).([]domain.Ticket), args.Error(1)
}

func (m *MockTicketService) UpdateTicketStatusByTicketID(ctx context.Context, id uuid.UUID, status domain.TicketStatus) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func setupTicketRouterWithAuth(h *TicketHandler, companyID uuid.UUID) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("companyID", companyID)
		c.Next()
	})
	r.POST("/api/v1/tickets", h.CreateTicket)
	r.POST("/api/v1/tickets/bulk", h.CreateTickets)
	r.GET("/api/v1/tickets", h.GetTicketsByCompanyID)
	r.PATCH("/api/v1/tickets/:id", h.UpdateTicketStatus)
	return r
}

func TestCreateTicket_Success(t *testing.T) {
	companyID := uuid.New()
	depID := uuid.New()
	formID := uuid.New()
	svc := new(MockTicketService)
	svc.On("CreateTicket", mock.Anything, mock.Anything).Return(nil)

	h := NewTicketHandler(svc, &config.Config{})
	w := makeJSONRequest(
		setupTicketRouterWithAuth(h, companyID),
		http.MethodPost,
		"/api/v1/tickets",
		map[string]any{
			"title":          "Test Ticket",
			"message":        "some message",
			"description":    "desc",
			"predict_status": "failed",
			"priority":       "high",
			"department_id":  depID,
			"form_id":        formID,
		},
	)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "create ticket success")
	svc.AssertExpectations(t)
}

func TestCreateTicket_InvalidBody(t *testing.T) {
	companyID := uuid.New()
	svc := new(MockTicketService)
	h := NewTicketHandler(svc, &config.Config{})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/tickets", bytes.NewBufferString("bad-json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	setupTicketRouterWithAuth(h, companyID).ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	svc.AssertNotCalled(t, "CreateTicket")
}

func TestCreateTicket_ServiceError(t *testing.T) {
	companyID := uuid.New()
	formID := uuid.New()
	svc := new(MockTicketService)
	svc.On("CreateTicket", mock.Anything, mock.Anything).Return(errors.New("db error"))

	h := NewTicketHandler(svc, &config.Config{})
	w := makeJSONRequest(
		setupTicketRouterWithAuth(h, companyID),
		http.MethodPost,
		"/api/v1/tickets",
		map[string]any{
			"title":          "T1",
			"message":        "M1",
			"description":    "D1",
			"predict_status": "failed",
			"priority":       "high",
			"form_id":        formID.String(),
		},
	)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	svc.AssertExpectations(t)
}

func TestCreateTickets_Success(t *testing.T) {
	companyID := uuid.New()
	formID := uuid.New()
	svc := new(MockTicketService)
	svc.On("CreateTickets", mock.Anything, mock.Anything).Return(nil)

	h := NewTicketHandler(svc, &config.Config{})
	w := makeJSONRequest(
		setupTicketRouterWithAuth(h, companyID),
		http.MethodPost,
		"/api/v1/tickets/bulk",
		[]map[string]any{
			{
				"title":          "T1",
				"message":        "M1",
				"description":    "D1",
				"predict_status": "failed",
				"priority":       "high",
				"form_id":        formID.String(),
			},
			{
				"title":          "T1",
				"message":        "M1",
				"description":    "D1",
				"predict_status": "failed",
				"priority":       "high",
				"form_id":        formID.String(),
			},
		},
	)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "create ticket success")
	svc.AssertExpectations(t)
}

func TestCreateTickets_InvalidBody(t *testing.T) {
	companyID := uuid.New()
	svc := new(MockTicketService)
	h := NewTicketHandler(svc, &config.Config{})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/tickets/bulk", bytes.NewBufferString("bad-json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	setupTicketRouterWithAuth(h, companyID).ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	svc.AssertNotCalled(t, "CreateTickets")
}

func TestCreateTickets_ServiceError(t *testing.T) {
	companyID := uuid.New()
	formID := uuid.New()
	svc := new(MockTicketService)
	svc.On("CreateTickets", mock.Anything, mock.Anything).Return(errors.New("db error"))

	h := NewTicketHandler(svc, &config.Config{})
	w := makeJSONRequest(
		setupTicketRouterWithAuth(h, companyID),
		http.MethodPost,
		"/api/v1/tickets/bulk",
		[]map[string]any{{
			"title":          "T1",
			"message":        "M1",
			"description":    "D1",
			"predict_status": "failed",
			"priority":       "high",
			"form_id":        formID.String(),
		}},
	)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	svc.AssertExpectations(t)
}

func TestGetTicketsByCompanyID_Success(t *testing.T) {
	companyID := uuid.New()
	svc := new(MockTicketService)
	svc.On("GetTicketsByCompanyID", mock.Anything, companyID).
		Return([]domain.Ticket{{Title: "Bug Report"}}, nil)

	h := NewTicketHandler(svc, &config.Config{})
	w := makeJSONRequest(setupTicketRouterWithAuth(h, companyID), http.MethodGet, "/api/v1/tickets", nil)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Bug Report")
	svc.AssertExpectations(t)
}

func TestGetTicketsByCompanyID_MissingCompanyID(t *testing.T) {
	svc := new(MockTicketService)
	h := NewTicketHandler(svc, &config.Config{})

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/api/v1/tickets", h.GetTicketsByCompanyID)

	w := makeJSONRequest(r, http.MethodGet, "/api/v1/tickets", nil)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	svc.AssertNotCalled(t, "GetTicketsByCompanyID")
}

func TestGetTicketsByCompanyID_ServiceError(t *testing.T) {
	companyID := uuid.New()
	svc := new(MockTicketService)
	svc.On("GetTicketsByCompanyID", mock.Anything, companyID).
		Return([]domain.Ticket{}, errors.New("db error"))

	h := NewTicketHandler(svc, &config.Config{})
	w := makeJSONRequest(setupTicketRouterWithAuth(h, companyID), http.MethodGet, "/api/v1/tickets", nil)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	svc.AssertExpectations(t)
}

func TestUpdateTicketStatus_Success(t *testing.T) {
	companyID := uuid.New()
	ticketID := uuid.New()
	svc := new(MockTicketService)
	svc.On("UpdateTicketStatusByTicketID", mock.Anything, ticketID, domain.TicketStatus("opened")).Return(nil)

	h := NewTicketHandler(svc, &config.Config{})
	w := makeJSONRequest(
		setupTicketRouterWithAuth(h, companyID),
		http.MethodPatch,
		"/api/v1/tickets/"+ticketID.String(),
		map[string]any{"status": "opened"},
	)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "ticket status updated")
	svc.AssertExpectations(t)
}

func TestUpdateTicketStatus_InvalidUUID(t *testing.T) {
	companyID := uuid.New()
	svc := new(MockTicketService)
	h := NewTicketHandler(svc, &config.Config{})

	w := makeJSONRequest(
		setupTicketRouterWithAuth(h, companyID),
		http.MethodPatch,
		"/api/v1/tickets/not-a-uuid",
		map[string]any{"status": "resolved"},
	)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	svc.AssertNotCalled(t, "UpdateTicketStatusByTicketID")
}

func TestUpdateTicketStatus_InvalidBody(t *testing.T) {
	companyID := uuid.New()
	ticketID := uuid.New()
	svc := new(MockTicketService)
	h := NewTicketHandler(svc, &config.Config{})

	req := httptest.NewRequest(http.MethodPatch, "/api/v1/tickets/"+ticketID.String(), bytes.NewBufferString("bad-json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	setupTicketRouterWithAuth(h, companyID).ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	svc.AssertNotCalled(t, "UpdateTicketStatusByTicketID")
}

func TestUpdateTicketStatus_ServiceError(t *testing.T) {
	companyID := uuid.New()
	ticketID := uuid.New()
	svc := new(MockTicketService)
	svc.On("UpdateTicketStatusByTicketID", mock.Anything, ticketID, mock.Anything).Return(errors.New("db error"))

	h := NewTicketHandler(svc, &config.Config{})
	w := makeJSONRequest(
		setupTicketRouterWithAuth(h, companyID),
		http.MethodPatch,
		"/api/v1/tickets/"+ticketID.String(),
		map[string]any{"status": "opened"},
	)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	svc.AssertExpectations(t)
}
