package service

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/auhmaugmaufm/predict-ticket-department-backend/internal/dto"
	"github.com/stretchr/testify/assert"
)

func TestSendFormsToAI_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/predict-LLM", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)
		assert.NotEmpty(t, r.Header.Get("X-HMAC-Signature"))

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(dto.AIResponse{ /* ใส่ field ที่มีใน struct */ })
	}))
	defer server.Close()

	svc := NewAIService(server.URL, "test-hmac-key")
	result, err := svc.SendFormsToAI(context.Background(), []dto.CompanyFormItems{
		{ /* ใส่ข้อมูล mock */ },
	})

	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestSendFormsToAI_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	svc := NewAIService(server.URL, "test-hmac-key")
	result, err := svc.SendFormsToAI(context.Background(), []dto.CompanyFormItems{})

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "500")
}

func TestSendFormsToAI_InvalidURL(t *testing.T) {
	svc := NewAIService("http://invalid-host-that-does-not-exist", "key")
	result, err := svc.SendFormsToAI(context.Background(), []dto.CompanyFormItems{})

	assert.Error(t, err)
	assert.Nil(t, result)
}
