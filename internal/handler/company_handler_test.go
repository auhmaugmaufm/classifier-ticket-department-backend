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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCompanyService struct {
	mock.Mock
}

func (m *MockCompanyService) Login(ctx context.Context, email, password string) (string, error) {
	args := m.Called(ctx, email, password)
	return args.String(0), args.Error(1)
}

func (m *MockCompanyService) Register(ctx context.Context, email, password string) error {
	args := m.Called(ctx, email, password)
	return args.Error(0)
}

func setupRouter(h *CompanyHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/api/v1/register", h.Register)
	r.POST("/api/v1/login", h.Login)
	return r
}

func makeRequest(router *gin.Engine, method, path string, body any) *httptest.ResponseRecorder {
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(method, path, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func TestRegister_Success(t *testing.T) {
	svc := new(MockCompanyService)
	svc.On("Register", mock.Anything, "test@example.com", "password123").Return(nil)

	h := NewCompanyHandler(svc, &config.Config{})
	w := makeRequest(setupRouter(h), http.MethodPost, "/api/v1/register", map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	})

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "register company success")
	svc.AssertExpectations(t)
}

func TestRegister_EmailAlreadyExists(t *testing.T) {
	svc := new(MockCompanyService)
	svc.On("Register", mock.Anything, "dup@example.com", "password123").
		Return(domain.ErrEmailAlreadyExists)

	h := NewCompanyHandler(svc, &config.Config{})
	w := makeRequest(setupRouter(h), http.MethodPost, "/api/v1/register", map[string]string{
		"email":    "dup@example.com",
		"password": "password123",
	})

	assert.Equal(t, http.StatusConflict, w.Code)
	svc.AssertExpectations(t)
}

func TestRegister_InternalServerError(t *testing.T) {
	svc := new(MockCompanyService)
	svc.On("Register", mock.Anything, "test@example.com", "password123").
		Return(errors.New("db connection failed"))

	h := NewCompanyHandler(svc, &config.Config{})
	w := makeRequest(setupRouter(h), http.MethodPost, "/api/v1/register", map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	})

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	svc.AssertExpectations(t)
}

func TestRegister_InvalidBody(t *testing.T) {
	svc := new(MockCompanyService)
	h := NewCompanyHandler(svc, &config.Config{})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/register", bytes.NewBufferString("not-json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	setupRouter(h).ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	svc.AssertNotCalled(t, "Register")
}

func TestLogin_Success(t *testing.T) {
	svc := new(MockCompanyService)
	svc.On("Login", mock.Anything, "test@example.com", "password123").Return("jwt-token-abc", nil)

	h := NewCompanyHandler(svc, &config.Config{})
	w := makeRequest(setupRouter(h), http.MethodPost, "/api/v1/login", map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	})

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "jwt-token-abc")
	svc.AssertExpectations(t)
}

func TestLogin_InvalidCredentials(t *testing.T) {
	svc := new(MockCompanyService)
	svc.On("Login", mock.Anything, "test@example.com", "wrongpass").
		Return("", domain.ErrInvalidCredentials)

	h := NewCompanyHandler(svc, &config.Config{})
	w := makeRequest(setupRouter(h), http.MethodPost, "/api/v1/login", map[string]string{
		"email":    "test@example.com",
		"password": "wrongpass",
	})

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	svc.AssertExpectations(t)
}

func TestLogin_InternalServerError(t *testing.T) {
	svc := new(MockCompanyService)
	svc.On("Login", mock.Anything, "test@example.com", "password123").
		Return("", errors.New("unexpected error"))

	h := NewCompanyHandler(svc, &config.Config{})
	w := makeRequest(setupRouter(h), http.MethodPost, "/api/v1/login", map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	})

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	svc.AssertExpectations(t)
}

func TestLogin_InvalidBody(t *testing.T) {
	svc := new(MockCompanyService)
	h := NewCompanyHandler(svc, &config.Config{})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewBufferString("{bad json}"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	setupRouter(h).ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	svc.AssertNotCalled(t, "Login")
}
