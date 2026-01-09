package destination

import (
	"context"
	"sort"
	"time"
)

// Repository define a interface para acesso a dados
type Repository interface {
	GetCountryMetadata(ctx context.Context, countryCode string) (*CountryMetadata, error)
	GetAllCountries(ctx context.Context) ([]CountryMetadata, error)
	GetMarketDataByNCM(ctx context.Context, ncm string, year, month int) ([]MarketData, error)
	GetMarketDataByNCMAndCountry(ctx context.Context, ncm, countryCode string, year, month int) (*MarketData, error)
}

// ServiceInterface define a interface do serviço de destinos
type ServiceInterface interface {
	RecommendDestinations(ctx context.Context, req SimulatorRequest) (*SimulatorResponse, error)
}

// Service implementa a lógica de negócio do simulador
type Service struct {
	repo    Repository
	weights ScoringWeights
}

// NewService cria uma nova instância do Service
func NewService(repo Repository) *Service {
	return &Service{
		repo:    repo,
		weights: DefaultScoringWeights(),
	}
}

// RecommendDestinations gera recomendações de destinos de exportação
func (s *Service) RecommendDestinations(ctx context.Context, req SimulatorRequest) (*SimulatorResponse, error) {
	startTime := time.Now()

	// Valida request
	if err := req.ValidateSimulatorRequest(); err != nil {
		return nil, err
	}

	// Busca dados de mercado do NCM (últimos 3 meses)
	now := time.Now()
	marketData, err := s.repo.GetMarketDataByNCM(ctx, req.NCM, now.Year(), int(now.Month()))
	if err != nil {
		return nil, err
	}

	if len(marketData) == 0 {
		return nil, ErrNoDataAvailable
	}

	// Busca metadados dos países
	countries, err := s.repo.GetAllCountries(ctx)
	if err != nil {
		return nil, err
	}

	// Cria map de países para lookup rápido
	countryMap := make(map[string]CountryMetadata)
	for _, c := range countries {
		countryMap[c.Code] = c
	}

	// Gera recomendações
	recommendations := make([]DestinationRecommendation, 0)
	maxMarketSize := 0.0
	maxGrowthRate := 0.0
	maxPrice := 0.0
	maxDistance := 0.0

	// Primeiro pass: calcular máximos para normalização
	for _, data := range marketData {
		if data.TotalValueUSD > maxMarketSize {
			maxMarketSize = data.TotalValueUSD
		}
		if data.GrowthRatePct > maxGrowthRate {
			maxGrowthRate = data.GrowthRatePct
		}
		if data.AvgPricePerKgUSD > maxPrice {
			maxPrice = data.AvgPricePerKgUSD
		}

		if country, exists := countryMap[data.CountryCode]; exists {
			if float64(country.DistanceBrazilKm) > maxDistance {
				maxDistance = float64(country.DistanceBrazilKm)
			}
		}
	}

	// Segundo pass: criar recomendações com scores
	for _, data := range marketData {
		country, exists := countryMap[data.CountryCode]
		if !exists {
			continue
		}

		// Filtra por países específicos (se fornecido)
		if len(req.Countries) > 0 {
			found := false
			for _, c := range req.Countries {
				if c == country.Code {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		rec := DestinationRecommendation{
			CountryCode:      country.Code,
			CountryName:      country.NamePt,
			MarketSizeUSD:    int64(data.TotalValueUSD),
			GrowthRatePct:    data.GrowthRatePct,
			PricePerKgUSD:    data.AvgPricePerKgUSD,
			DistanceKm:       country.DistanceBrazilKm,
			Region:           country.Region,
			FlagEmoji:        country.FlagEmoji,
			EstimatedMarginPct: s.estimateMargin(data.AvgPricePerKgUSD),
			LogisticsCostUSD: s.estimateLogisticsCost(country.DistanceBrazilKm, req.VolumeKg),
			TariffRatePct:    s.estimateTariff(country.Region),
			LeadTimeDays:     s.estimateLeadTime(country.DistanceBrazilKm),
		}

		// Calcula score
		rec.Score = rec.CalculateScore(s.weights, maxMarketSize, maxGrowthRate, maxPrice, maxDistance)
		rec.Demand = rec.GetDemandLevel()
		rec.RecommendationReason = rec.GetRecommendationReason(s.weights)

		recommendations = append(recommendations, rec)
	}

	// Ordena por score (maior primeiro)
	sort.Slice(recommendations, func(i, j int) bool {
		return recommendations[i].Score > recommendations[j].Score
	})

	// Adiciona ranking
	for i := range recommendations {
		recommendations[i].Rank = i + 1
	}

	// Limita resultados
	if len(recommendations) > req.MaxResults {
		recommendations = recommendations[:req.MaxResults]
	}

	processingTime := time.Since(startTime).Milliseconds()

	return &SimulatorResponse{
		Destinations: recommendations,
		Metadata: SimulatorMetadata{
			NCM:              req.NCM,
			ProductName:      "", // TODO: Buscar nome do produto
			AnalysisDate:     time.Now(),
			TotalDestinations: len(recommendations),
			ProcessingTimeMs: processingTime,
		},
	}, nil
}

// estimateMargin estima margem baseada no preço
func (s *Service) estimateMargin(pricePerKg float64) float64 {
	// Lógica simplificada: margens maiores para preços mais altos
	if pricePerKg > 50 {
		return 35.0
	}
	if pricePerKg > 20 {
		return 25.0
	}
	return 15.0
}

// estimateLogisticsCost estima custo logístico
func (s *Service) estimateLogisticsCost(distanceKm int, volumeKg *float64) float64 {
	// Custo base por km
	costPerKm := 0.05

	volume := 1000.0 // default 1 tonelada
	if volumeKg != nil {
		volume = *volumeKg
	}

	// Custo = distância × custo/km × volume (com economia de escala)
	scaleFactor := 1.0
	if volume > 10000 {
		scaleFactor = 0.8 // 20% desconto para volumes > 10 toneladas
	}

	return float64(distanceKm) * costPerKm * (volume / 1000.0) * scaleFactor
}

// estimateTariff estima tarifa baseada na região
func (s *Service) estimateTariff(region string) float64 {
	switch region {
	case "Americas":
		return 8.0 // Mercosul tem tarifas menores
	case "Europe":
		return 12.0
	case "Asia":
		return 15.0
	case "Africa":
		return 18.0
	default:
		return 12.0
	}
}

// estimateLeadTime estima tempo de entrega
func (s *Service) estimateLeadTime(distanceKm int) int {
	// Fórmula simplificada: ~500km por dia de transporte
	days := distanceKm / 500

	// Adiciona tempo de processamento aduaneiro
	processingDays := 7

	return days + processingDays
}
