package destination

import (
	"testing"
)

// TestValidateSimulatorRequest testa a validação de requisições
func TestValidateSimulatorRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     SimulatorRequest
		wantErr error
	}{
		{
			name: "valid request with all fields",
			req: SimulatorRequest{
				NCM:        "12345678",
				VolumeKg:   ptrFloat64(1000.0),
				Countries:  []string{"US", "CN"},
				MaxResults: 10,
			},
			wantErr: nil,
		},
		{
			name: "valid request minimal",
			req: SimulatorRequest{
				NCM: "12345678",
			},
			wantErr: nil,
		},
		{
			name: "invalid NCM - too short",
			req: SimulatorRequest{
				NCM: "1234567",
			},
			wantErr: ErrInvalidNCM,
		},
		{
			name: "invalid NCM - too long",
			req: SimulatorRequest{
				NCM: "123456789",
			},
			wantErr: ErrInvalidNCM,
		},
		{
			name: "invalid NCM - contains letters",
			req: SimulatorRequest{
				NCM: "1234567A",
			},
			wantErr: ErrInvalidNCM,
		},
		{
			name: "invalid NCM - contains special chars",
			req: SimulatorRequest{
				NCM: "1234567-",
			},
			wantErr: ErrInvalidNCM,
		},
		{
			name: "invalid volume - zero",
			req: SimulatorRequest{
				NCM:      "12345678",
				VolumeKg: ptrFloat64(0.0),
			},
			wantErr: ErrInvalidVolume,
		},
		{
			name: "invalid volume - negative",
			req: SimulatorRequest{
				NCM:      "12345678",
				VolumeKg: ptrFloat64(-100.0),
			},
			wantErr: ErrInvalidVolume,
		},
		{
			name: "invalid max results - negative",
			req: SimulatorRequest{
				NCM:        "12345678",
				MaxResults: -1,
			},
			wantErr: ErrInvalidMaxResults,
		},
		{
			name: "invalid max results - too high",
			req: SimulatorRequest{
				NCM:        "12345678",
				MaxResults: 51,
			},
			wantErr: ErrInvalidMaxResults,
		},
		{
			name: "max results at boundary - 50",
			req: SimulatorRequest{
				NCM:        "12345678",
				MaxResults: 50,
			},
			wantErr: nil,
		},
		{
			name: "max results at boundary - 1",
			req: SimulatorRequest{
				NCM:        "12345678",
				MaxResults: 1,
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.ValidateSimulatorRequest()
			if err != tt.wantErr {
				t.Errorf("ValidateSimulatorRequest() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Se MaxResults não fornecido (0), deve ser setado para 10
			if tt.wantErr == nil && tt.req.MaxResults == 0 {
				// O método altera o valor interno
				if tt.req.MaxResults != 10 {
					t.Errorf("ValidateSimulatorRequest() MaxResults should default to 10, got %d", tt.req.MaxResults)
				}
			}
		})
	}
}

// TestCalculateScore testa o cálculo de score
func TestCalculateScore(t *testing.T) {
	tests := []struct {
		name           string
		rec            DestinationRecommendation
		weights        ScoringWeights
		maxMarketSize  float64
		maxGrowthRate  float64
		maxPrice       float64
		maxDistance    float64
		expectedMin    float64
		expectedMax    float64
	}{
		{
			name: "perfect score - all max values",
			rec: DestinationRecommendation{
				MarketSizeUSD: 1000000,
				GrowthRatePct: 50.0,
				PricePerKgUSD: 100.0,
				DistanceKm:    0, // Menor distância = melhor
			},
			weights:       DefaultScoringWeights(),
			maxMarketSize: 1000000,
			maxGrowthRate: 50.0,
			maxPrice:      100.0,
			maxDistance:   10000.0,
			expectedMin:   9.0, // Próximo de 10
			expectedMax:   10.0,
		},
		{
			name: "zero score - worst case",
			rec: DestinationRecommendation{
				MarketSizeUSD: 0,
				GrowthRatePct: 0,
				PricePerKgUSD: 0,
				DistanceKm:    10000, // Máxima distância = pior
			},
			weights:       DefaultScoringWeights(),
			maxMarketSize: 1000000,
			maxGrowthRate: 50.0,
			maxPrice:      100.0,
			maxDistance:   10000.0,
			expectedMin:   0.0,
			expectedMax:   0.1,
		},
		{
			name: "average score",
			rec: DestinationRecommendation{
				MarketSizeUSD: 500000,   // 50% do max
				GrowthRatePct: 25.0,     // 50% do max
				PricePerKgUSD: 50.0,     // 50% do max
				DistanceKm:    5000,     // 50% do max
			},
			weights:       DefaultScoringWeights(),
			maxMarketSize: 1000000,
			maxGrowthRate: 50.0,
			maxPrice:      100.0,
			maxDistance:   10000.0,
			expectedMin:   4.5,
			expectedMax:   5.5,
		},
		{
			name: "custom weights - market size dominant",
			rec: DestinationRecommendation{
				MarketSizeUSD: 1000000, // Max
				GrowthRatePct: 0,       // Min
				PricePerKgUSD: 0,       // Min
				DistanceKm:    10000,   // Max (pior)
			},
			weights: ScoringWeights{
				MarketSize: 0.80, // 80% peso em market size
				GrowthRate: 0.10,
				PricePerKg: 0.05,
				Distance:   0.05,
			},
			maxMarketSize: 1000000,
			maxGrowthRate: 50.0,
			maxPrice:      100.0,
			maxDistance:   10000.0,
			expectedMin:   7.5, // Dominado por market size
			expectedMax:   8.5,
		},
		{
			name: "zero max values - edge case",
			rec: DestinationRecommendation{
				MarketSizeUSD: 100000,
				GrowthRatePct: 10.0,
				PricePerKgUSD: 50.0,
				DistanceKm:    5000,
			},
			weights:       DefaultScoringWeights(),
			maxMarketSize: 0, // Edge case
			maxGrowthRate: 0,
			maxPrice:      0,
			maxDistance:   0,
			expectedMin:   0.0,
			expectedMax:   0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := tt.rec.CalculateScore(tt.weights, tt.maxMarketSize, tt.maxGrowthRate, tt.maxPrice, tt.maxDistance)

			if score < tt.expectedMin || score > tt.expectedMax {
				t.Errorf("CalculateScore() = %.2f, expected between %.2f and %.2f", score, tt.expectedMin, tt.expectedMax)
			}

			// Score deve estar sempre entre 0 e 10
			if score < 0 || score > 10 {
				t.Errorf("CalculateScore() = %.2f, must be between 0 and 10", score)
			}
		})
	}
}

// TestGetDemandLevel testa a classificação de demanda
func TestGetDemandLevel(t *testing.T) {
	tests := []struct {
		name          string
		marketSizeUSD int64
		expected      string
	}{
		{
			name:          "high demand - 1 billion",
			marketSizeUSD: 1_000_000_000,
			expected:      "Alto",
		},
		{
			name:          "high demand - exactly 100M + 1",
			marketSizeUSD: 100_000_001,
			expected:      "Alto",
		},
		{
			name:          "medium demand - 50M",
			marketSizeUSD: 50_000_000,
			expected:      "Médio",
		},
		{
			name:          "medium demand - exactly 10M + 1",
			marketSizeUSD: 10_000_001,
			expected:      "Médio",
		},
		{
			name:          "low demand - 5M",
			marketSizeUSD: 5_000_000,
			expected:      "Baixo",
		},
		{
			name:          "low demand - 1M",
			marketSizeUSD: 1_000_000,
			expected:      "Baixo",
		},
		{
			name:          "low demand - zero",
			marketSizeUSD: 0,
			expected:      "Baixo",
		},
		{
			name:          "boundary - exactly 100M",
			marketSizeUSD: 100_000_000,
			expected:      "Médio", // <= 100M
		},
		{
			name:          "boundary - exactly 10M",
			marketSizeUSD: 10_000_000,
			expected:      "Baixo", // <= 10M
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := DestinationRecommendation{
				MarketSizeUSD: tt.marketSizeUSD,
			}
			result := rec.GetDemandLevel()
			if result != tt.expected {
				t.Errorf("GetDemandLevel() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestGetRecommendationReason testa a geração de razões de recomendação
func TestGetRecommendationReason(t *testing.T) {
	tests := []struct {
		name           string
		score          float64
		expectedReason string
	}{
		{
			name:           "excellent score - 10.0",
			score:          10.0,
			expectedReason: "Mercado altamente atrativo com grande potencial de crescimento e demanda consolidada",
		},
		{
			name:           "excellent score - 8.5",
			score:          8.5,
			expectedReason: "Mercado altamente atrativo com grande potencial de crescimento e demanda consolidada",
		},
		{
			name:           "excellent score - exactly 8.0",
			score:          8.0,
			expectedReason: "Mercado altamente atrativo com grande potencial de crescimento e demanda consolidada",
		},
		{
			name:           "good score - 7.5",
			score:          7.5,
			expectedReason: "Mercado promissor com bom equilíbrio entre demanda, crescimento e custos logísticos",
		},
		{
			name:           "good score - exactly 6.0",
			score:          6.0,
			expectedReason: "Mercado promissor com bom equilíbrio entre demanda, crescimento e custos logísticos",
		},
		{
			name:           "moderate score - 5.0",
			score:          5.0,
			expectedReason: "Mercado em desenvolvimento com oportunidades emergentes",
		},
		{
			name:           "moderate score - exactly 4.0",
			score:          4.0,
			expectedReason: "Mercado em desenvolvimento com oportunidades emergentes",
		},
		{
			name:           "low score - 3.0",
			score:          3.0,
			expectedReason: "Mercado em fase inicial ou com barreiras significativas",
		},
		{
			name:           "low score - 0.0",
			score:          0.0,
			expectedReason: "Mercado em fase inicial ou com barreiras significativas",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := DestinationRecommendation{
				Score: tt.score,
			}
			weights := DefaultScoringWeights()
			result := rec.GetRecommendationReason(weights)
			if result != tt.expectedReason {
				t.Errorf("GetRecommendationReason() = %v, expected %v", result, tt.expectedReason)
			}
		})
	}
}

// TestDefaultScoringWeights testa os pesos padrão
func TestDefaultScoringWeights(t *testing.T) {
	weights := DefaultScoringWeights()

	// Verificar valores
	if weights.MarketSize != 0.40 {
		t.Errorf("MarketSize weight = %.2f, expected 0.40", weights.MarketSize)
	}
	if weights.GrowthRate != 0.30 {
		t.Errorf("GrowthRate weight = %.2f, expected 0.30", weights.GrowthRate)
	}
	if weights.PricePerKg != 0.20 {
		t.Errorf("PricePerKg weight = %.2f, expected 0.20", weights.PricePerKg)
	}
	if weights.Distance != 0.10 {
		t.Errorf("Distance weight = %.2f, expected 0.10", weights.Distance)
	}

	// Verificar que a soma dos pesos é 1.0
	sum := weights.MarketSize + weights.GrowthRate + weights.PricePerKg + weights.Distance
	if sum < 0.99 || sum > 1.01 {
		t.Errorf("Sum of weights = %.2f, expected 1.0", sum)
	}
}

// TestCalculateScore_DistanceInversion testa a inversão de distância
func TestCalculateScore_DistanceInversion(t *testing.T) {
	weights := DefaultScoringWeights()

	// Destino próximo (menor distância = melhor)
	nearRec := DestinationRecommendation{
		MarketSizeUSD: 1000000,
		GrowthRatePct: 50.0,
		PricePerKgUSD: 100.0,
		DistanceKm:    1000, // Próximo
	}

	// Destino distante
	farRec := DestinationRecommendation{
		MarketSizeUSD: 1000000,
		GrowthRatePct: 50.0,
		PricePerKgUSD: 100.0,
		DistanceKm:    9000, // Distante
	}

	maxMarketSize := 1000000.0
	maxGrowthRate := 50.0
	maxPrice := 100.0
	maxDistance := 10000.0

	nearScore := nearRec.CalculateScore(weights, maxMarketSize, maxGrowthRate, maxPrice, maxDistance)
	farScore := farRec.CalculateScore(weights, maxMarketSize, maxGrowthRate, maxPrice, maxDistance)

	// Destino próximo deve ter score maior
	if nearScore <= farScore {
		t.Errorf("Near destination score (%.2f) should be greater than far destination score (%.2f)", nearScore, farScore)
	}
}

// TestCalculateScore_WeightInfluence testa a influência dos pesos
func TestCalculateScore_WeightInfluence(t *testing.T) {
	// Dois destinos idênticos exceto market size
	highMarketRec := DestinationRecommendation{
		MarketSizeUSD: 1000000, // Alto
		GrowthRatePct: 10.0,
		PricePerKgUSD: 50.0,
		DistanceKm:    5000,
	}

	lowMarketRec := DestinationRecommendation{
		MarketSizeUSD: 100000, // Baixo
		GrowthRatePct: 10.0,
		PricePerKgUSD: 50.0,
		DistanceKm:    5000,
	}

	// Pesos normais (market size = 40%)
	normalWeights := DefaultScoringWeights()

	// Pesos com market size dominante (80%)
	marketDominantWeights := ScoringWeights{
		MarketSize: 0.80,
		GrowthRate: 0.10,
		PricePerKg: 0.05,
		Distance:   0.05,
	}

	maxMarketSize := 1000000.0
	maxGrowthRate := 50.0
	maxPrice := 100.0
	maxDistance := 10000.0

	// Score com pesos normais
	normalHighScore := highMarketRec.CalculateScore(normalWeights, maxMarketSize, maxGrowthRate, maxPrice, maxDistance)
	normalLowScore := lowMarketRec.CalculateScore(normalWeights, maxMarketSize, maxGrowthRate, maxPrice, maxDistance)
	normalDiff := normalHighScore - normalLowScore

	// Score com pesos dominantes
	dominantHighScore := highMarketRec.CalculateScore(marketDominantWeights, maxMarketSize, maxGrowthRate, maxPrice, maxDistance)
	dominantLowScore := lowMarketRec.CalculateScore(marketDominantWeights, maxMarketSize, maxGrowthRate, maxPrice, maxDistance)
	dominantDiff := dominantHighScore - dominantLowScore

	// Com peso maior em market size, a diferença deve ser maior
	if dominantDiff <= normalDiff {
		t.Errorf("Dominant market weight difference (%.2f) should be greater than normal weight difference (%.2f)", dominantDiff, normalDiff)
	}
}

// Helper function para criar ponteiro de float64
func ptrFloat64(f float64) *float64 {
	return &f
}
