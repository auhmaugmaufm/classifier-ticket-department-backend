package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/auhmaugmaufm/predict-ticket-department-backend/internal/auth"
	"github.com/auhmaugmaufm/predict-ticket-department-backend/internal/dto"
)

type AIService struct {
	client  *http.Client
	baseURL string
	hmacKey string
}

func NewAIService(baseURL string, hmacKey string) *AIService {
	return &AIService{
		client:  &http.Client{Timeout: 30 * time.Second},
		baseURL: baseURL,
		hmacKey: hmacKey,
	}
}

func (s *AIService) SendFormsToAI(ctx context.Context, data []dto.CompanyFormItems) (*dto.AIResponse, error) {
	body, err := json.Marshal(dto.AIDataRequest{Data: data})
	fmt.Println(string(body))
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.baseURL+"/api/v1/predict-LLM", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	message := string(body)
	signature := auth.ComputeHMAC(s.hmacKey, message)
	req.Header.Set("X-HMAC-Signature", "sha256="+signature)

	fmt.Printf("%+v\n", req)

	res, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("ai_service error status: %d", res.StatusCode)
	}

	var result dto.AIResponse
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("ai_service decode: %w", err)
	}

	return &result, nil

}
