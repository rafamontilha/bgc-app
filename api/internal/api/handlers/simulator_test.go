package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"bgc-app/internal/business/destination"

	"github.com/gin-gonic/gin"
)

// MockService implementa destination.Service para testes
type MockDestinationService struct {
	RecommendDestinationsFunc func(ctx context.Context, req destination.SimulatorRequest) (*destination.SimulatorResponse, error)
}

func (m *MockDestinationService) RecommendDestinations(ctx context.Context, req destination.SimulatorRequest) (*destination.SimulatorResponse, error) {
	if m.RecommendDestinationsFunc != nil {
		return m.RecommendDestinationsFunc(ctx, req)
	}
	return nil, errors.New("not implemented")
}

// setupTestRouter cria um router Gin para testes
func setupTestRouter(handler *SimulatorHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/v1/simulator/destinations", handler.SimulateDestinations)
	return r
}

// TestSimulateDestinations_Success testa o fluxo de sucesso
func TestSimulateDestinations_Success(t *testing.T) {
	mockService := &MockDestinationService{
		RecommendDestinationsFunc: func(ctx context.Context, req destination.SimulatorRequest) (*destination.SimulatorResponse, error) {
			return &destination.SimulatorResponse{
				Destinations: []destination.DestinationRecommendation{
					{
						CountryCode:          "US",
						CountryName:          "Estados Unidos",
						Score:                8.5,
						Rank:                 1,
						Demand:               "Alto",
						EstimatedMarginPct:   25.0,
						LogisticsCostUSD:     500.0,
						TariffRatePct:        8.0,
						LeadTimeDays:         20,
						MarketSizeUSD:        1000000,
						GrowthRatePct:        15.5,
						PricePerKgUSD:        100.0,
						DistanceKm:           7500,
						Region:               "Americas",
						FlagEmoji:            "🇺🇸",
						RecommendationReason: "Mercado altamente atrativo",
					},
				},
				Metadata: destination.SimulatorMetadata{
					NCM:               "12345678",
					ProductName:       "Test Product",
					AnalysisDate:      time.Now(),
					TotalDestinations: 1,
					ProcessingTimeMs:  50,
				},
			}, nil
		},
	}

	// Criar handler adaptado
	handler := &SimulatorHandler{
		service: mockService,
	}
	router := setupTestRouter(handler)

	// Criar request
	reqBody := map[string]interface{}{
		"ncm":         "12345678",
		"volume_kg":   1000.0,
		"max_results": 10,
	}
	jsonBody, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/v1/simulator/destinations", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Verificar status code
	if w.Code != http.StatusOK {
		t.Errorf("Status code = %d, expected %d", w.Code, http.StatusOK)
	}

	// Verificar response body
	var resp destination.SimulatorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(resp.Destinations) != 1 {
		t.Errorf("Response has %d destinations, expected 1", len(resp.Destinations))
	}

	if resp.Destinations[0].CountryCode != "US" {
		t.Errorf("Country code = %s, expected US", resp.Destinations[0].CountryCode)
	}

	// Verificar que processing time foi atualizado no handler (pode ser 0 se muito rápido)
	// O importante é que foi setado pelo handler, não pelo mock
	if resp.Metadata.ProcessingTimeMs < 0 {
		t.Error("ProcessingTimeMs should not be negative")
	}
}

// TestSimulateDestinations_InvalidJSON testa JSON inválido
func TestSimulateDestinations_InvalidJSON(t *testing.T) {
	handler := &SimulatorHandler{
		service: &MockDestinationService{},
	}
	router := setupTestRouter(handler)

	// JSON malformado
	req := httptest.NewRequest(http.MethodPost, "/v1/simulator/destinations", bytes.NewBufferString("{invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Deve retornar 400 Bad Request
	if w.Code != http.StatusBadRequest {
		t.Errorf("Status code = %d, expected %d", w.Code, http.StatusBadRequest)
	}

	var errResp ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &errResp); err != nil {
		t.Fatalf("Failed to unmarshal error response: %v", err)
	}

	if errResp.Error != "invalid_request" {
		t.Errorf("Error code = %s, expected invalid_request", errResp.Error)
	}
}

// TestSimulateDestinations_MissingNCM testa NCM faltando
func TestSimulateDestinations_MissingNCM(t *testing.T) {
	handler := &SimulatorHandler{
		service: &MockDestinationService{},
	}
	router := setupTestRouter(handler)

	// Request sem NCM
	reqBody := map[string]interface{}{
		"volume_kg": 1000.0,
	}
	jsonBody, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/v1/simulator/destinations", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Status code = %d, expected %d", w.Code, http.StatusBadRequest)
	}
}

// TestSimulateDestinations_ValidationError testa erro de validação
func TestSimulateDestinations_ValidationError(t *testing.T) {
	mockService := &MockDestinationService{
		RecommendDestinationsFunc: func(ctx context.Context, req destination.SimulatorRequest) (*destination.SimulatorResponse, error) {
			return nil, destination.ErrInvalidNCM
		},
	}

	handler := &SimulatorHandler{
		service: mockService,
	}
	router := setupTestRouter(handler)

	reqBody := map[string]interface{}{
		"ncm": "12345678",
	}
	jsonBody, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/v1/simulator/destinations", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Status code = %d, expected %d", w.Code, http.StatusBadRequest)
	}

	var errResp ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &errResp); err != nil {
		t.Fatalf("Failed to unmarshal error response: %v", err)
	}

	if errResp.Error != "validation_error" {
		t.Errorf("Error code = %s, expected validation_error", errResp.Error)
	}
}

// TestSimulateDestinations_NCMNotFound testa NCM não encontrado
func TestSimulateDestinations_NCMNotFound(t *testing.T) {
	mockService := &MockDestinationService{
		RecommendDestinationsFunc: func(ctx context.Context, req destination.SimulatorRequest) (*destination.SimulatorResponse, error) {
			return nil, destination.ErrNCMNotFound
		},
	}

	handler := &SimulatorHandler{
		service: mockService,
	}
	router := setupTestRouter(handler)

	reqBody := map[string]interface{}{
		"ncm": "99999999",
	}
	jsonBody, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/v1/simulator/destinations", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Status code = %d, expected %d", w.Code, http.StatusNotFound)
	}

	var errResp ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &errResp); err != nil {
		t.Fatalf("Failed to unmarshal error response: %v", err)
	}

	if errResp.Error != "ncm_not_found" {
		t.Errorf("Error code = %s, expected ncm_not_found", errResp.Error)
	}
}

// TestSimulateDestinations_NoDataAvailable testa dados não disponíveis
func TestSimulateDestinations_NoDataAvailable(t *testing.T) {
	mockService := &MockDestinationService{
		RecommendDestinationsFunc: func(ctx context.Context, req destination.SimulatorRequest) (*destination.SimulatorResponse, error) {
			return nil, destination.ErrNoDataAvailable
		},
	}

	handler := &SimulatorHandler{
		service: mockService,
	}
	router := setupTestRouter(handler)

	reqBody := map[string]interface{}{
		"ncm": "12345678",
	}
	jsonBody, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/v1/simulator/destinations", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Status code = %d, expected %d", w.Code, http.StatusNotFound)
	}

	var errResp ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &errResp); err != nil {
		t.Fatalf("Failed to unmarshal error response: %v", err)
	}

	if errResp.Error != "no_data_available" {
		t.Errorf("Error code = %s, expected no_data_available", errResp.Error)
	}
}

// TestSimulateDestinations_InsufficientData testa dados insuficientes
func TestSimulateDestinations_InsufficientData(t *testing.T) {
	mockService := &MockDestinationService{
		RecommendDestinationsFunc: func(ctx context.Context, req destination.SimulatorRequest) (*destination.SimulatorResponse, error) {
			return nil, destination.ErrInsufficientData
		},
	}

	handler := &SimulatorHandler{
		service: mockService,
	}
	router := setupTestRouter(handler)

	reqBody := map[string]interface{}{
		"ncm": "12345678",
	}
	jsonBody, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/v1/simulator/destinations", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("Status code = %d, expected %d", w.Code, http.StatusUnprocessableEntity)
	}

	var errResp ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &errResp); err != nil {
		t.Fatalf("Failed to unmarshal error response: %v", err)
	}

	if errResp.Error != "insufficient_data" {
		t.Errorf("Error code = %s, expected insufficient_data", errResp.Error)
	}
}

// TestSimulateDestinations_InternalError testa erro interno
func TestSimulateDestinations_InternalError(t *testing.T) {
	mockService := &MockDestinationService{
		RecommendDestinationsFunc: func(ctx context.Context, req destination.SimulatorRequest) (*destination.SimulatorResponse, error) {
			return nil, errors.New("unexpected database error")
		},
	}

	handler := &SimulatorHandler{
		service: mockService,
	}
	router := setupTestRouter(handler)

	reqBody := map[string]interface{}{
		"ncm": "12345678",
	}
	jsonBody, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/v1/simulator/destinations", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Status code = %d, expected %d", w.Code, http.StatusInternalServerError)
	}

	var errResp ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &errResp); err != nil {
		t.Fatalf("Failed to unmarshal error response: %v", err)
	}

	if errResp.Error != "internal_error" {
		t.Errorf("Error code = %s, expected internal_error", errResp.Error)
	}

	// Error details deve estar presente
	if errResp.Details == "" {
		t.Error("Error details should not be empty for internal errors")
	}
}

// TestSimulateDestinations_AllValidationErrors testa todos os erros de validação
func TestSimulateDestinations_AllValidationErrors(t *testing.T) {
	tests := []struct {
		name           string
		serviceError   error
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "invalid NCM",
			serviceError:   destination.ErrInvalidNCM,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "validation_error",
		},
		{
			name:           "invalid volume",
			serviceError:   destination.ErrInvalidVolume,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "validation_error",
		},
		{
			name:           "invalid max results",
			serviceError:   destination.ErrInvalidMaxResults,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "validation_error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockDestinationService{
				RecommendDestinationsFunc: func(ctx context.Context, req destination.SimulatorRequest) (*destination.SimulatorResponse, error) {
					return nil, tt.serviceError
				},
			}

			handler := &SimulatorHandler{
				service: mockService,
			}
			router := setupTestRouter(handler)

			reqBody := map[string]interface{}{
				"ncm": "12345678",
			}
			jsonBody, _ := json.Marshal(reqBody)

			req := httptest.NewRequest(http.MethodPost, "/v1/simulator/destinations", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Status code = %d, expected %d", w.Code, tt.expectedStatus)
			}

			var errResp ErrorResponse
			if err := json.Unmarshal(w.Body.Bytes(), &errResp); err != nil {
				t.Fatalf("Failed to unmarshal error response: %v", err)
			}

			if errResp.Error != tt.expectedError {
				t.Errorf("Error code = %s, expected %s", errResp.Error, tt.expectedError)
			}
		})
	}
}

// TestSimulateDestinations_WithCountries testa request com filtro de países
func TestSimulateDestinations_WithCountries(t *testing.T) {
	var capturedRequest destination.SimulatorRequest

	mockService := &MockDestinationService{
		RecommendDestinationsFunc: func(ctx context.Context, req destination.SimulatorRequest) (*destination.SimulatorResponse, error) {
			capturedRequest = req
			return &destination.SimulatorResponse{
				Destinations: []destination.DestinationRecommendation{},
				Metadata: destination.SimulatorMetadata{
					NCM:               req.NCM,
					TotalDestinations: 0,
				},
			}, nil
		},
	}

	handler := &SimulatorHandler{
		service: mockService,
	}
	router := setupTestRouter(handler)

	reqBody := map[string]interface{}{
		"ncm":       "12345678",
		"countries": []string{"US", "CN", "DE"},
	}
	jsonBody, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/v1/simulator/destinations", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status code = %d, expected %d", w.Code, http.StatusOK)
	}

	// Verificar que countries foi capturado
	if len(capturedRequest.Countries) != 3 {
		t.Errorf("Countries length = %d, expected 3", len(capturedRequest.Countries))
	}

	expectedCountries := map[string]bool{"US": true, "CN": true, "DE": true}
	for _, country := range capturedRequest.Countries {
		if !expectedCountries[country] {
			t.Errorf("Unexpected country: %s", country)
		}
	}
}

// TestSimulateDestinations_ProcessingTimeAdded testa que processing time é adicionado
func TestSimulateDestinations_ProcessingTimeAdded(t *testing.T) {
	mockService := &MockDestinationService{
		RecommendDestinationsFunc: func(ctx context.Context, req destination.SimulatorRequest) (*destination.SimulatorResponse, error) {
			// Simular processamento
			time.Sleep(10 * time.Millisecond)
			return &destination.SimulatorResponse{
				Destinations: []destination.DestinationRecommendation{},
				Metadata: destination.SimulatorMetadata{
					NCM:              req.NCM,
					ProcessingTimeMs: 0, // Handler deve sobrescrever
				},
			}, nil
		},
	}

	handler := &SimulatorHandler{
		service: mockService,
	}
	router := setupTestRouter(handler)

	reqBody := map[string]interface{}{
		"ncm": "12345678",
	}
	jsonBody, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/v1/simulator/destinations", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	var resp destination.SimulatorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Processing time deve ser >= 10ms
	if resp.Metadata.ProcessingTimeMs < 10 {
		t.Errorf("ProcessingTimeMs = %d, expected >= 10", resp.Metadata.ProcessingTimeMs)
	}
}

// TestNewSimulatorHandler testa criação do handler
func TestNewSimulatorHandler(t *testing.T) {
	mockService := &MockDestinationService{}
	handler := NewSimulatorHandler(mockService)

	if handler == nil {
		t.Error("NewSimulatorHandler() should not return nil")
	}

	if handler.service != mockService {
		t.Error("NewSimulatorHandler() should store the service")
	}
}
