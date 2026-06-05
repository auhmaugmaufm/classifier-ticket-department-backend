package handler

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/auhmaugmaufm/predict-ticket-department-backend/internal/domain"
	"github.com/auhmaugmaufm/predict-ticket-department-backend/pkg/config"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockFormService struct {
	mock.Mock
}

func (m *MockFormService) SubmitForm(ctx context.Context, form *domain.Form) error {
	args := m.Called(ctx, form)
	return args.Error(0)
}

func (m *MockFormService) GetSubmitFormByCompanyID(ctx context.Context, companyID uuid.UUID) ([]domain.Form, error) {
	args := m.Called(ctx, companyID)
	return args.Get(0).([]domain.Form), args.Error(1)
}

func (m *MockFormService) GetSubmitFormPerDayByCompanyID(ctx context.Context, companyID uuid.UUID, dateStr string) ([]domain.Form, error) {
	args := m.Called(ctx, companyID, dateStr)
	return args.Get(0).([]domain.Form), args.Error(1)
}

func setupFormRouterWithAuth(h *FormHandler, companyID uuid.UUID) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("companyID", companyID)
		c.Next()
	})
	r.POST("/forms/submit", h.SubmitForm)
	r.GET("/api/v1/forms", h.GetSubmitFormCompanyID)
	r.GET("/api/v1/forms/per-day", h.GetSubmitFormPerDayByCompanyID)
	return r
}

func TestSubmitForm_Success(t *testing.T) {
	companyID := uuid.New()
	svc := new(MockFormService)
	svc.On("SubmitForm", mock.Anything, mock.Anything).Return(nil)

	h := NewFormHandler(svc, &config.Config{})
	w := makeJSONRequest(
		setupFormRouterWithAuth(h, companyID),
		http.MethodPost,
		"/forms/submit",
		map[string]any{
			"company_id":  companyID,
			"title":       "Support Request",
			"description": "Need help",
		},
	)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "submit form success")
	svc.AssertExpectations(t)
}

func TestSubmitForm_InvalidBody(t *testing.T) {
	companyID := uuid.New()
	svc := new(MockFormService)
	h := NewFormHandler(svc, &config.Config{})

	req := httptest.NewRequest(http.MethodPost, "/forms/submit", bytes.NewBufferString("bad-json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	setupFormRouterWithAuth(h, companyID).ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	svc.AssertNotCalled(t, "SubmitForm")
}

func TestSubmitForm_ServiceError(t *testing.T) {
	companyID := uuid.New()
	svc := new(MockFormService)
	svc.On("SubmitForm", mock.Anything, mock.Anything).Return(errors.New("db error"))

	h := NewFormHandler(svc, &config.Config{})
	w := makeJSONRequest(
		setupFormRouterWithAuth(h, companyID),
		http.MethodPost,
		"/forms/submit",
		map[string]any{"company_id": companyID, "title": "T", "description": "D"},
	)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	svc.AssertExpectations(t)
}

func TestGetSubmitFormCompanyID_Success(t *testing.T) {
	companyID := uuid.New()
	svc := new(MockFormService)
	svc.On("GetSubmitFormByCompanyID", mock.Anything, companyID).
		Return([]domain.Form{{Title: "My Form"}}, nil)

	h := NewFormHandler(svc, &config.Config{})
	w := makeJSONRequest(setupFormRouterWithAuth(h, companyID), http.MethodGet, "/api/v1/forms", nil)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "My Form")
	svc.AssertExpectations(t)
}

func TestGetSubmitFormCompanyID_MissingCompanyID(t *testing.T) {
	svc := new(MockFormService)
	h := NewFormHandler(svc, &config.Config{})

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/api/v1/forms", h.GetSubmitFormCompanyID)

	w := makeJSONRequest(r, http.MethodGet, "/api/v1/forms", nil)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	svc.AssertNotCalled(t, "GetSubmitFormByCompanyID")
}

func TestGetSubmitFormCompanyID_ServiceError(t *testing.T) {
	companyID := uuid.New()
	svc := new(MockFormService)
	svc.On("GetSubmitFormByCompanyID", mock.Anything, companyID).
		Return([]domain.Form{}, errors.New("db error"))

	h := NewFormHandler(svc, &config.Config{})
	w := makeJSONRequest(setupFormRouterWithAuth(h, companyID), http.MethodGet, "/api/v1/forms", nil)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	svc.AssertExpectations(t)
}

func TestGetSubmitFormPerDay_Success(t *testing.T) {
	companyID := uuid.New()
	date := "2025-01-15"
	svc := new(MockFormService)
	svc.On("GetSubmitFormPerDayByCompanyID", mock.Anything, companyID, date).
		Return([]domain.Form{{Title: "Daily Form"}}, nil)

	h := NewFormHandler(svc, &config.Config{})
	w := makeJSONRequest(
		setupFormRouterWithAuth(h, companyID),
		http.MethodGet,
		"/api/v1/forms/per-day?date="+date,
		nil,
	)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Daily Form")
	svc.AssertExpectations(t)
}

func TestGetSubmitFormPerDay_DefaultYesterday(t *testing.T) {
	companyID := uuid.New()
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	svc := new(MockFormService)
	svc.On("GetSubmitFormPerDayByCompanyID", mock.Anything, companyID, yesterday).
		Return([]domain.Form{}, nil)

	h := NewFormHandler(svc, &config.Config{})
	// ไม่ส่ง ?date= เลย ต้อง default เป็น yesterday
	w := makeJSONRequest(setupFormRouterWithAuth(h, companyID), http.MethodGet, "/api/v1/forms/per-day", nil)

	assert.Equal(t, http.StatusOK, w.Code)
	svc.AssertExpectations(t)
}

func TestGetSubmitFormPerDay_InvalidDateFormat(t *testing.T) {
	companyID := uuid.New()
	svc := new(MockFormService)
	h := NewFormHandler(svc, &config.Config{})

	w := makeJSONRequest(
		setupFormRouterWithAuth(h, companyID),
		http.MethodGet,
		"/api/v1/forms/per-day?date=01-15-2025", // format ผิด
		nil,
	)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid date format")
	svc.AssertNotCalled(t, "GetSubmitFormPerDayByCompanyID")
}

func TestGetSubmitFormPerDay_MissingCompanyID(t *testing.T) {
	svc := new(MockFormService)
	h := NewFormHandler(svc, &config.Config{})

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/api/v1/forms/per-day", h.GetSubmitFormPerDayByCompanyID)

	w := makeJSONRequest(r, http.MethodGet, "/api/v1/forms/per-day?date=2025-01-15", nil)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	svc.AssertNotCalled(t, "GetSubmitFormPerDayByCompanyID")
}

func TestGetSubmitFormPerDay_ServiceError(t *testing.T) {
	companyID := uuid.New()
	date := "2025-01-15"
	svc := new(MockFormService)
	svc.On("GetSubmitFormPerDayByCompanyID", mock.Anything, companyID, date).
		Return([]domain.Form{}, errors.New("db error"))

	h := NewFormHandler(svc, &config.Config{})
	w := makeJSONRequest(
		setupFormRouterWithAuth(h, companyID),
		http.MethodGet,
		"/api/v1/forms/per-day?date="+date,
		nil,
	)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	svc.AssertExpectations(t)
}
