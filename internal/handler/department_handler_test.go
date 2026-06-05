package handler

import (
	"bytes"
	"context"
	"encoding/json"
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

type MockDepartmentService struct {
	mock.Mock
}

func (m *MockDepartmentService) AddDepartments(ctx context.Context, departments []domain.Department) error {
	args := m.Called(ctx, departments)
	return args.Error(0)
}

func (m *MockDepartmentService) GetDepartmentsByCompanyID(ctx context.Context, companyID uuid.UUID) ([]domain.Department, error) {
	args := m.Called(ctx, companyID)
	return args.Get(0).([]domain.Department), args.Error(1)
}

func (m *MockDepartmentService) UpdateDepartmentStatus(ctx context.Context, id uuid.UUID, isActive bool) error {
	args := m.Called(ctx, id, isActive)
	return args.Error(0)
}

func (m *MockDepartmentService) DeleteDepartmentByID(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func makeJSONRequest(router *gin.Engine, method, path string, body any) *httptest.ResponseRecorder {
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(method, path, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func setupDepartmentRouterWithAuth(h *DepartmentHandler, companyID uuid.UUID) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("companyID", companyID)
		c.Next()
	})
	r.POST("/api/v1/departments", h.CreateDepartment)
	r.GET("/api/v1/departments", h.GetDepartmentsFromToken)
	r.PATCH("/api/v1/departments/:id", h.UpdateDepartmentStatus)
	r.DELETE("/api/v1/departments/:id", h.DeleteDepartmentByID)
	return r
}

func TestCreateDepartment_Success(t *testing.T) {
	companyID := uuid.New()
	svc := new(MockDepartmentService)

	expectedDeps := []domain.Department{
		{DepartmentName: "IT", CompanyID: companyID},
		{DepartmentName: "HR", CompanyID: companyID},
	}
	svc.On("AddDepartments", mock.Anything, expectedDeps).Return(nil)

	h := NewDepartmentHandler(svc, &config.Config{})
	w := makeJSONRequest(
		setupDepartmentRouterWithAuth(h, companyID),
		http.MethodPost,
		"/api/v1/departments",
		map[string]any{"department_name": []string{"IT", "HR"}},
	)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "add departments success")
	svc.AssertExpectations(t)
}

func TestCreateDepartment_InvalidBody(t *testing.T) {
	companyID := uuid.New()
	svc := new(MockDepartmentService)

	h := NewDepartmentHandler(svc, &config.Config{})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/departments", bytes.NewBufferString("bad-json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	setupDepartmentRouterWithAuth(h, companyID).ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	svc.AssertNotCalled(t, "AddDepartments")
}

func TestCreateDepartment_MissingCompanyID(t *testing.T) {
	svc := new(MockDepartmentService)
	h := NewDepartmentHandler(svc, &config.Config{})

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/api/v1/departments", h.CreateDepartment)

	w := makeJSONRequest(r, http.MethodPost, "/api/v1/departments",
		map[string]any{"department_name": []string{"IT"}},
	)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	svc.AssertNotCalled(t, "AddDepartments")
}

func TestCreateDepartment_ServiceError(t *testing.T) {
	companyID := uuid.New()
	svc := new(MockDepartmentService)
	svc.On("AddDepartments", mock.Anything, mock.Anything).Return(errors.New("db error"))

	h := NewDepartmentHandler(svc, &config.Config{})
	w := makeJSONRequest(
		setupDepartmentRouterWithAuth(h, companyID),
		http.MethodPost,
		"/api/v1/departments",
		map[string]any{"department_name": []string{"IT"}},
	)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	svc.AssertExpectations(t)
}

func TestGetDepartmentsByCompanyID_Success(t *testing.T) {
	companyID := uuid.New()
	svc := new(MockDepartmentService)
	mockDeps := []domain.Department{
		{DepartmentName: "IT", CompanyID: companyID},
	}
	svc.On("GetDepartmentsByCompanyID", mock.Anything, companyID).Return(mockDeps, nil)

	h := NewDepartmentHandler(svc, &config.Config{})

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/api/v1/internal/departments/:company_id", h.GetDepartmentsByCompanyID)

	w := makeJSONRequest(r, http.MethodGet, "/api/v1/internal/departments/"+companyID.String(), nil)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "IT")
	svc.AssertExpectations(t)
}

func TestGetDepartmentsByCompanyID_InvalidUUID(t *testing.T) {
	svc := new(MockDepartmentService)
	h := NewDepartmentHandler(svc, &config.Config{})

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/api/v1/internal/departments/:company_id", h.GetDepartmentsByCompanyID)

	w := makeJSONRequest(r, http.MethodGet, "/api/v1/internal/departments/not-a-uuid", nil)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	svc.AssertNotCalled(t, "GetDepartmentsByCompanyID")
}

func TestGetDepartmentsByCompanyID_ServiceError(t *testing.T) {
	companyID := uuid.New()
	svc := new(MockDepartmentService)
	svc.On("GetDepartmentsByCompanyID", mock.Anything, companyID).
		Return([]domain.Department{}, errors.New("db error"))

	h := NewDepartmentHandler(svc, &config.Config{})

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/api/v1/internal/departments/:company_id", h.GetDepartmentsByCompanyID)

	w := makeJSONRequest(r, http.MethodGet, "/api/v1/internal/departments/"+companyID.String(), nil)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	svc.AssertExpectations(t)
}

func TestGetDepartmentsFromToken_Success(t *testing.T) {
	companyID := uuid.New()
	svc := new(MockDepartmentService)
	mockDeps := []domain.Department{
		{DepartmentName: "Finance", CompanyID: companyID},
	}
	svc.On("GetDepartmentsByCompanyID", mock.Anything, companyID).Return(mockDeps, nil)

	h := NewDepartmentHandler(svc, &config.Config{})
	w := makeJSONRequest(
		setupDepartmentRouterWithAuth(h, companyID),
		http.MethodGet,
		"/api/v1/departments",
		nil,
	)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Finance")
	svc.AssertExpectations(t)
}

func TestGetDepartmentsFromToken_MissingCompanyID(t *testing.T) {
	svc := new(MockDepartmentService)
	h := NewDepartmentHandler(svc, &config.Config{})

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/api/v1/departments", h.GetDepartmentsFromToken)

	w := makeJSONRequest(r, http.MethodGet, "/api/v1/departments", nil)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	svc.AssertNotCalled(t, "GetDepartmentsByCompanyID")
}

func TestGetDepartmentsFromToken_ServiceError(t *testing.T) {
	companyID := uuid.New()
	svc := new(MockDepartmentService)
	svc.On("GetDepartmentsByCompanyID", mock.Anything, companyID).
		Return([]domain.Department{}, errors.New("db error"))

	h := NewDepartmentHandler(svc, &config.Config{})
	w := makeJSONRequest(
		setupDepartmentRouterWithAuth(h, companyID),
		http.MethodGet,
		"/api/v1/departments",
		nil,
	)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	svc.AssertExpectations(t)
}

func TestUpdateDepartmentStatus_Success(t *testing.T) {
	companyID := uuid.New()
	depID := uuid.New()
	svc := new(MockDepartmentService)
	svc.On("UpdateDepartmentStatus", mock.Anything, depID, true).Return(nil)

	h := NewDepartmentHandler(svc, &config.Config{})
	w := makeJSONRequest(
		setupDepartmentRouterWithAuth(h, companyID),
		http.MethodPatch,
		"/api/v1/departments/"+depID.String(),
		map[string]any{"is_active": true},
	)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "department status updated")
	svc.AssertExpectations(t)
}

func TestUpdateDepartmentStatus_InvalidUUID(t *testing.T) {
	companyID := uuid.New()
	svc := new(MockDepartmentService)
	h := NewDepartmentHandler(svc, &config.Config{})

	w := makeJSONRequest(
		setupDepartmentRouterWithAuth(h, companyID),
		http.MethodPatch,
		"/api/v1/departments/not-a-uuid",
		map[string]any{"is_active": true},
	)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	svc.AssertNotCalled(t, "UpdateDepartmentStatus")
}

func TestUpdateDepartmentStatus_InvalidBody(t *testing.T) {
	companyID := uuid.New()
	depID := uuid.New()
	svc := new(MockDepartmentService)
	h := NewDepartmentHandler(svc, &config.Config{})

	req := httptest.NewRequest(http.MethodPatch, "/api/v1/departments/"+depID.String(), bytes.NewBufferString("bad-json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	setupDepartmentRouterWithAuth(h, companyID).ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	svc.AssertNotCalled(t, "UpdateDepartmentStatus")
}

func TestUpdateDepartmentStatus_ServiceError(t *testing.T) {
	companyID := uuid.New()
	depID := uuid.New()
	svc := new(MockDepartmentService)
	svc.On("UpdateDepartmentStatus", mock.Anything, depID, false).Return(errors.New("db error"))

	h := NewDepartmentHandler(svc, &config.Config{})
	w := makeJSONRequest(
		setupDepartmentRouterWithAuth(h, companyID),
		http.MethodPatch,
		"/api/v1/departments/"+depID.String(),
		map[string]any{"is_active": false},
	)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	svc.AssertExpectations(t)
}

func TestDeleteDepartmentByID_Success(t *testing.T) {
	companyID := uuid.New()
	depID := uuid.New()
	svc := new(MockDepartmentService)
	svc.On("DeleteDepartmentByID", mock.Anything, depID).Return(nil)

	h := NewDepartmentHandler(svc, &config.Config{})
	w := makeJSONRequest(
		setupDepartmentRouterWithAuth(h, companyID),
		http.MethodDelete,
		"/api/v1/departments/"+depID.String(),
		nil,
	)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "department deleted successfully")
	svc.AssertExpectations(t)
}

func TestDeleteDepartmentByID_InvalidUUID(t *testing.T) {
	companyID := uuid.New()
	svc := new(MockDepartmentService)
	h := NewDepartmentHandler(svc, &config.Config{})

	w := makeJSONRequest(
		setupDepartmentRouterWithAuth(h, companyID),
		http.MethodDelete,
		"/api/v1/departments/not-a-uuid",
		nil,
	)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	svc.AssertNotCalled(t, "DeleteDepartmentByID")
}

func TestDeleteDepartmentByID_ServiceError(t *testing.T) {
	companyID := uuid.New()
	depID := uuid.New()
	svc := new(MockDepartmentService)
	svc.On("DeleteDepartmentByID", mock.Anything, depID).Return(errors.New("db error"))

	h := NewDepartmentHandler(svc, &config.Config{})
	w := makeJSONRequest(
		setupDepartmentRouterWithAuth(h, companyID),
		http.MethodDelete,
		"/api/v1/departments/"+depID.String(),
		nil,
	)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	svc.AssertExpectations(t)
}
