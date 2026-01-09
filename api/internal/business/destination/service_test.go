package destination

import (
	"context"
	"errors"
	"testing"
	"time"
)

// MockRepository implementa a interface Repository para testes
type MockRepository struct {
	GetCountryMetadataFunc           func(ctx context.Context, countryCode string) (*CountryMetadata, error)
	GetAllCountriesFunc              func(ctx context.Context) ([]CountryMetadata, error)
	GetMarketDataByNCMFunc           func(ctx context.Context, ncm string, year, month int) ([]MarketData, error)
	GetMarketDataByNCMAndCountryFunc func(ctx context.Context, ncm, countryCode string, year, month int) (*MarketData, error)
}

func (m *MockRepository) GetCountryMetadata(ctx context.Context, countryCode string) (*CountryMetadata, error) {
	if m.GetCountryMetadataFunc != nil {
		return m.GetCountryMetadataFunc(ctx, countryCode)
	}
	return nil, errors.New("not implemented")
}

func (m *MockRepository) GetAllCountries(ctx context.Context) ([]CountryMetadata, error) {
	if m.GetAllCountriesFunc != nil {
		return m.GetAllCountriesFunc(ctx)
	}
	return nil, errors.New("not implemented")
}

func (m *MockRepository) GetMarketDataByNCM(ctx context.Context, ncm string, year, month int) ([]MarketData, error) {
	if m.GetMarketDataByNCMFunc != nil {
		return m.GetMarketDataByNCMFunc(ctx, ncm, year, month)
	}
	return nil, errors.New("not implemented")
}

func (m *MockRepository) GetMarketDataByNCMAndCountry(ctx context.Context, ncm, countryCode string, year, month int) (*MarketData, error) {
	if m.GetMarketDataByNCMAndCountryFunc != nil {
		return m.GetMarketDataByNCMAndCountryFunc(ctx, ncm, countryCode, year, month)
	}
	return nil, errors.New("not implemented")
}

// TestNewService testa a criação de um novo service
func TestNewService(t *testing.T) {
	mockRepo := &MockRepository{}
	service := NewService(mockRepo)

	if service == nil {
		t.Error("NewService() should not return nil")
	}

	if service.repo != mockRepo {
		t.Error("NewService() should store the repository")
	}

	// Verificar que os pesos padrão são usados
	expectedWeights := DefaultScoringWeights()
	if service.weights != expectedWeights {
		t.Errorf("NewService() weights = %+v, expected %+v", service.weights, expectedWeights)
	}
}

// TestRecommendDestinations_Success testa o fluxo completo de sucesso
func TestRecommendDestinations_Success(t *testing.T) {
	mockRepo := &MockRepository{
		GetMarketDataByNCMFunc: func(ctx context.Context, ncm string, year, month int) ([]MarketData, error) {
			return []MarketData{
				{
					NCM:              "12345678",
					CountryCode:      "US",
					Year:             2024,
					Month:            12,
					TotalValueUSD:    1000000,
					TotalWeightKg:    10000,
					AvgPricePerKgUSD: 100,
					TransactionCount: 50,
					GrowthRatePct:    15.5,
				},
				{
					NCM:              "12345678",
					CountryCode:      "CN",
					Year:             2024,
					Month:            12,
					TotalValueUSD:    800000,
					TotalWeightKg:    8000,
					AvgPricePerKgUSD: 100,
					TransactionCount: 40,
					GrowthRatePct:    20.0,
				},
			}, nil
		},
		GetAllCountriesFunc: func(ctx context.Context) ([]CountryMetadata, error) {
			return []CountryMetadata{
				{
					Code:             "US",
					NamePt:           "Estados Unidos",
					Region:           "Americas",
					DistanceBrazilKm: 7500,
					FlagEmoji:        "🇺🇸",
				},
				{
					Code:             "CN",
					NamePt:           "China",
					Region:           "Asia",
					DistanceBrazilKm: 17000,
					FlagEmoji:        "🇨🇳",
				},
			}, nil
		},
	}

	service := NewService(mockRepo)
	req := SimulatorRequest{
		NCM:        "12345678",
		MaxResults: 10,
	}

	resp, err := service.RecommendDestinations(context.Background(), req)

	if err != nil {
		t.Fatalf("RecommendDestinations() error = %v, expected nil", err)
	}

	if resp == nil {
		t.Fatal("RecommendDestinations() returned nil response")
	}

	if len(resp.Destinations) != 2 {
		t.Errorf("RecommendDestinations() returned %d destinations, expected 2", len(resp.Destinations))
	}

	// Verificar que o primeiro item tem rank 1
	if resp.Destinations[0].Rank != 1 {
		t.Errorf("First destination rank = %d, expected 1", resp.Destinations[0].Rank)
	}

	// Verificar que o segundo item tem rank 2
	if resp.Destinations[1].Rank != 2 {
		t.Errorf("Second destination rank = %d, expected 2", resp.Destinations[1].Rank)
	}

	// Verificar que estão ordenados por score (decrescente)
	if resp.Destinations[0].Score < resp.Destinations[1].Score {
		t.Error("Destinations should be sorted by score in descending order")
	}

	// Verificar metadata
	if resp.Metadata.NCM != "12345678" {
		t.Errorf("Metadata NCM = %s, expected 12345678", resp.Metadata.NCM)
	}

	if resp.Metadata.TotalDestinations != 2 {
		t.Errorf("Metadata TotalDestinations = %d, expected 2", resp.Metadata.TotalDestinations)
	}

	if resp.Metadata.ProcessingTimeMs <= 0 {
		t.Error("ProcessingTimeMs should be greater than 0")
	}
}

// TestRecommendDestinations_InvalidRequest testa validação de request inválido
func TestRecommendDestinations_InvalidRequest(t *testing.T) {
	mockRepo := &MockRepository{}
	service := NewService(mockRepo)

	tests := []struct {
		name    string
		req     SimulatorRequest
		wantErr error
	}{
		{
			name: "invalid NCM",
			req: SimulatorRequest{
				NCM: "123", // Muito curto
			},
			wantErr: ErrInvalidNCM,
		},
		{
			name: "invalid volume",
			req: SimulatorRequest{
				NCM:      "12345678",
				VolumeKg: ptrFloat64(-100),
			},
			wantErr: ErrInvalidVolume,
		},
		{
			name: "invalid max results",
			req: SimulatorRequest{
				NCM:        "12345678",
				MaxResults: 100, // Muito alto
			},
			wantErr: ErrInvalidMaxResults,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.RecommendDestinations(context.Background(), tt.req)
			if err != tt.wantErr {
				t.Errorf("RecommendDestinations() error = %v, expected %v", err, tt.wantErr)
			}
		})
	}
}

// TestRecommendDestinations_NoDataAvailable testa quando não há dados
func TestRecommendDestinations_NoDataAvailable(t *testing.T) {
	mockRepo := &MockRepository{
		GetMarketDataByNCMFunc: func(ctx context.Context, ncm string, year, month int) ([]MarketData, error) {
			return []MarketData{}, nil // Sem dados
		},
	}

	service := NewService(mockRepo)
	req := SimulatorRequest{
		NCM:        "12345678",
		MaxResults: 10,
	}

	_, err := service.RecommendDestinations(context.Background(), req)

	if err != ErrNoDataAvailable {
		t.Errorf("RecommendDestinations() error = %v, expected ErrNoDataAvailable", err)
	}
}

// TestRecommendDestinations_RepositoryError testa erro no repositório
func TestRecommendDestinations_RepositoryError(t *testing.T) {
	expectedErr := errors.New("database connection failed")

	mockRepo := &MockRepository{
		GetMarketDataByNCMFunc: func(ctx context.Context, ncm string, year, month int) ([]MarketData, error) {
			return nil, expectedErr
		},
	}

	service := NewService(mockRepo)
	req := SimulatorRequest{
		NCM:        "12345678",
		MaxResults: 10,
	}

	_, err := service.RecommendDestinations(context.Background(), req)

	if err != expectedErr {
		t.Errorf("RecommendDestinations() error = %v, expected %v", err, expectedErr)
	}
}

// TestRecommendDestinations_CountryFilter testa filtro por países
func TestRecommendDestinations_CountryFilter(t *testing.T) {
	mockRepo := &MockRepository{
		GetMarketDataByNCMFunc: func(ctx context.Context, ncm string, year, month int) ([]MarketData, error) {
			return []MarketData{
				{CountryCode: "US", TotalValueUSD: 1000000, AvgPricePerKgUSD: 100, GrowthRatePct: 10},
				{CountryCode: "CN", TotalValueUSD: 800000, AvgPricePerKgUSD: 100, GrowthRatePct: 15},
				{CountryCode: "DE", TotalValueUSD: 600000, AvgPricePerKgUSD: 100, GrowthRatePct: 8},
			}, nil
		},
		GetAllCountriesFunc: func(ctx context.Context) ([]CountryMetadata, error) {
			return []CountryMetadata{
				{Code: "US", NamePt: "Estados Unidos", Region: "Americas", DistanceBrazilKm: 7500},
				{Code: "CN", NamePt: "China", Region: "Asia", DistanceBrazilKm: 17000},
				{Code: "DE", NamePt: "Alemanha", Region: "Europe", DistanceBrazilKm: 9500},
			}, nil
		},
	}

	service := NewService(mockRepo)
	req := SimulatorRequest{
		NCM:        "12345678",
		Countries:  []string{"US", "CN"}, // Filtrar apenas US e CN
		MaxResults: 10,
	}

	resp, err := service.RecommendDestinations(context.Background(), req)

	if err != nil {
		t.Fatalf("RecommendDestinations() error = %v, expected nil", err)
	}

	// Deve retornar apenas US e CN, não DE
	if len(resp.Destinations) != 2 {
		t.Errorf("RecommendDestinations() returned %d destinations, expected 2", len(resp.Destinations))
	}

	// Verificar que DE não está presente
	for _, dest := range resp.Destinations {
		if dest.CountryCode == "DE" {
			t.Error("Germany should be filtered out")
		}
	}
}

// TestRecommendDestinations_MaxResultsLimit testa limite de resultados
func TestRecommendDestinations_MaxResultsLimit(t *testing.T) {
	// Criar 20 países de teste
	marketData := make([]MarketData, 20)
	countries := make([]CountryMetadata, 20)

	for i := 0; i < 20; i++ {
		code := string(rune('A' + i))
		marketData[i] = MarketData{
			CountryCode:      code,
			TotalValueUSD:    float64(1000000 - i*10000),
			AvgPricePerKgUSD: 100,
			GrowthRatePct:    10,
		}
		countries[i] = CountryMetadata{
			Code:             code,
			NamePt:           "País " + code,
			Region:           "Americas",
			DistanceBrazilKm: 5000,
		}
	}

	mockRepo := &MockRepository{
		GetMarketDataByNCMFunc: func(ctx context.Context, ncm string, year, month int) ([]MarketData, error) {
			return marketData, nil
		},
		GetAllCountriesFunc: func(ctx context.Context) ([]CountryMetadata, error) {
			return countries, nil
		},
	}

	service := NewService(mockRepo)
	req := SimulatorRequest{
		NCM:        "12345678",
		MaxResults: 5, // Limitar a 5 resultados
	}

	resp, err := service.RecommendDestinations(context.Background(), req)

	if err != nil {
		t.Fatalf("RecommendDestinations() error = %v, expected nil", err)
	}

	if len(resp.Destinations) != 5 {
		t.Errorf("RecommendDestinations() returned %d destinations, expected 5", len(resp.Destinations))
	}
}

// TestEstimateMargin testa estimativa de margem
func TestEstimateMargin(t *testing.T) {
	service := NewService(&MockRepository{})

	tests := []struct {
		name         string
		pricePerKg   float64
		expectedMin  float64
		expectedMax  float64
	}{
		{
			name:        "high price - 100",
			pricePerKg:  100,
			expectedMin: 30,
			expectedMax: 40,
		},
		{
			name:        "medium price - 30",
			pricePerKg:  30,
			expectedMin: 20,
			expectedMax: 30,
		},
		{
			name:        "low price - 10",
			pricePerKg:  10,
			expectedMin: 10,
			expectedMax: 20,
		},
		{
			name:        "very low price - 1",
			pricePerKg:  1,
			expectedMin: 10,
			expectedMax: 20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			margin := service.estimateMargin(tt.pricePerKg)
			if margin < tt.expectedMin || margin > tt.expectedMax {
				t.Errorf("estimateMargin(%.2f) = %.2f, expected between %.2f and %.2f",
					tt.pricePerKg, margin, tt.expectedMin, tt.expectedMax)
			}
		})
	}
}

// TestEstimateLogisticsCost testa estimativa de custo logístico
func TestEstimateLogisticsCost(t *testing.T) {
	service := NewService(&MockRepository{})

	tests := []struct {
		name        string
		distanceKm  int
		volumeKg    *float64
		expectedMin float64
		expectedMax float64
	}{
		{
			name:        "short distance, small volume",
			distanceKm:  1000,
			volumeKg:    ptrFloat64(1000),
			expectedMin: 40,
			expectedMax: 60,
		},
		{
			name:        "long distance, small volume",
			distanceKm:  10000,
			volumeKg:    ptrFloat64(1000),
			expectedMin: 400,
			expectedMax: 600,
		},
		{
			name:        "short distance, large volume with discount",
			distanceKm:  1000,
			volumeKg:    ptrFloat64(20000), // > 10 toneladas
			expectedMin: 700,
			expectedMax: 900,
		},
		{
			name:        "nil volume uses default",
			distanceKm:  1000,
			volumeKg:    nil, // Deve usar 1000kg default
			expectedMin: 40,
			expectedMax: 60,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cost := service.estimateLogisticsCost(tt.distanceKm, tt.volumeKg)
			if cost < tt.expectedMin || cost > tt.expectedMax {
				t.Errorf("estimateLogisticsCost(%d, %.0f) = %.2f, expected between %.2f and %.2f",
					tt.distanceKm, derefFloat64(tt.volumeKg, 1000), cost, tt.expectedMin, tt.expectedMax)
			}
		})
	}
}

// TestEstimateTariff testa estimativa de tarifa por região
func TestEstimateTariff(t *testing.T) {
	service := NewService(&MockRepository{})

	tests := []struct {
		name     string
		region   string
		expected float64
	}{
		{"Americas - lowest", "Americas", 8.0},
		{"Europe - medium", "Europe", 12.0},
		{"Asia - higher", "Asia", 15.0},
		{"Africa - highest", "Africa", 18.0},
		{"Unknown - default", "Unknown", 12.0},
		{"Empty - default", "", 12.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tariff := service.estimateTariff(tt.region)
			if tariff != tt.expected {
				t.Errorf("estimateTariff(%s) = %.2f, expected %.2f", tt.region, tariff, tt.expected)
			}
		})
	}
}

// TestEstimateLeadTime testa estimativa de tempo de entrega
func TestEstimateLeadTime(t *testing.T) {
	service := NewService(&MockRepository{})

	tests := []struct {
		name        string
		distanceKm  int
		expectedMin int
		expectedMax int
	}{
		{
			name:        "very close - 500km",
			distanceKm:  500,
			expectedMin: 7,  // 1 dia transporte + 7 processamento
			expectedMax: 9,
		},
		{
			name:        "medium - 5000km",
			distanceKm:  5000,
			expectedMin: 15, // 10 dias transporte + 7 processamento
			expectedMax: 18,
		},
		{
			name:        "very far - 17000km (China)",
			distanceKm:  17000,
			expectedMin: 38, // 34 dias transporte + 7 processamento
			expectedMax: 45,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			leadTime := service.estimateLeadTime(tt.distanceKm)
			if leadTime < tt.expectedMin || leadTime > tt.expectedMax {
				t.Errorf("estimateLeadTime(%d) = %d, expected between %d and %d",
					tt.distanceKm, leadTime, tt.expectedMin, tt.expectedMax)
			}
		})
	}
}

// TestRecommendDestinations_Context testa cancelamento de contexto
func TestRecommendDestinations_Context(t *testing.T) {
	// Criar contexto já cancelado
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancela imediatamente

	mockRepo := &MockRepository{
		GetMarketDataByNCMFunc: func(ctx context.Context, ncm string, year, month int) ([]MarketData, error) {
			// Verificar se contexto foi cancelado
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
				return []MarketData{}, nil
			}
		},
	}

	service := NewService(mockRepo)
	req := SimulatorRequest{
		NCM:        "12345678",
		MaxResults: 10,
	}

	_, err := service.RecommendDestinations(ctx, req)

	if err != context.Canceled {
		t.Errorf("RecommendDestinations() error = %v, expected context.Canceled", err)
	}
}

// TestRecommendDestinations_ProcessingTime testa que processing time é calculado
func TestRecommendDestinations_ProcessingTime(t *testing.T) {
	mockRepo := &MockRepository{
		GetMarketDataByNCMFunc: func(ctx context.Context, ncm string, year, month int) ([]MarketData, error) {
			// Simular delay
			time.Sleep(10 * time.Millisecond)
			return []MarketData{
				{CountryCode: "US", TotalValueUSD: 1000000, AvgPricePerKgUSD: 100, GrowthRatePct: 10},
			}, nil
		},
		GetAllCountriesFunc: func(ctx context.Context) ([]CountryMetadata, error) {
			return []CountryMetadata{
				{Code: "US", NamePt: "Estados Unidos", Region: "Americas", DistanceBrazilKm: 7500},
			}, nil
		},
	}

	service := NewService(mockRepo)
	req := SimulatorRequest{
		NCM:        "12345678",
		MaxResults: 10,
	}

	resp, err := service.RecommendDestinations(context.Background(), req)

	if err != nil {
		t.Fatalf("RecommendDestinations() error = %v, expected nil", err)
	}

	// Processing time deve ser > 10ms devido ao sleep
	if resp.Metadata.ProcessingTimeMs < 10 {
		t.Errorf("ProcessingTimeMs = %d, expected >= 10", resp.Metadata.ProcessingTimeMs)
	}
}

// Helper para dereferenciar float64 pointer
func derefFloat64(ptr *float64, defaultVal float64) float64 {
	if ptr == nil {
		return defaultVal
	}
	return *ptr
}
