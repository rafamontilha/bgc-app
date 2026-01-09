package destination

import "time"

// DestinationRecommendation representa uma recomendação de destino de exportação
type DestinationRecommendation struct {
	CountryCode          string  `json:"country_code"`            // Código ISO do país (BR, US, CN)
	CountryName          string  `json:"country_name"`            // Nome do país
	Score                float64 `json:"score"`                   // Score geral (0-10)
	Rank                 int     `json:"rank"`                    // Ranking da recomendação
	Demand               string  `json:"demand"`                  // Alto, Médio, Baixo
	EstimatedMarginPct   float64 `json:"estimated_margin_pct"`    // Margem estimada (%)
	LogisticsCostUSD     float64 `json:"logistics_cost_usd"`      // Custo logístico estimado
	TariffRatePct        float64 `json:"tariff_rate_pct"`         // Taxa de tarifa de importação
	LeadTimeDays         int     `json:"lead_time_days"`          // Tempo estimado de entrega
	MarketSizeUSD        int64   `json:"market_size_usd"`         // Tamanho do mercado em USD
	GrowthRatePct        float64 `json:"growth_rate_pct"`         // Taxa de crescimento anual
	PricePerKgUSD        float64 `json:"price_per_kg_usd"`        // Preço médio por kg
	DistanceKm           int     `json:"distance_km"`             // Distância do Brasil
	Region               string  `json:"region"`                  // Região geográfica
	FlagEmoji            string  `json:"flag_emoji,omitempty"`    // Emoji da bandeira
	RecommendationReason string  `json:"recommendation_reason"`   // Explicação do score
}

// SimulatorRequest representa a requisição ao simulador
type SimulatorRequest struct {
	NCM         string   `json:"ncm" binding:"required,len=8"`         // NCM de 8 dígitos
	VolumeKg    *float64 `json:"volume_kg,omitempty"`                  // Volume em kg (opcional)
	Countries   []string `json:"countries,omitempty"`                  // Lista de países específicos (opcional)
	MaxResults  int      `json:"max_results,omitempty"`                // Número máximo de resultados (default: 10)
	IncludeAll  bool     `json:"include_all,omitempty"`                // Incluir todos os países disponíveis
}

// SimulatorResponse representa a resposta do simulador
type SimulatorResponse struct {
	Destinations []DestinationRecommendation `json:"destinations"` // Lista de destinos recomendados
	Metadata     SimulatorMetadata           `json:"metadata"`     // Metadados da análise
}

// SimulatorMetadata metadados sobre a análise
type SimulatorMetadata struct {
	NCM              string    `json:"ncm"`                // NCM analisado
	ProductName      string    `json:"product_name"`       // Nome do produto (se disponível)
	AnalysisDate     time.Time `json:"analysis_date"`      // Data da análise
	TotalDestinations int      `json:"total_destinations"` // Total de destinos analisados
	CacheHit         bool      `json:"cache_hit"`          // Foi cache hit?
	CacheLevel       string    `json:"cache_level,omitempty"` // Nível do cache (l1, l2, l3)
	ProcessingTimeMs int64     `json:"processing_time_ms"` // Tempo de processamento
}

// CountryMetadata representa os metadados de um país
type CountryMetadata struct {
	Code                   string   `json:"code"`
	NamePt                 string   `json:"name_pt"`
	NameEn                 string   `json:"name_en"`
	Region                 string   `json:"region"`
	Subregion              string   `json:"subregion"`
	GDPUSD                 *int64   `json:"gdp_usd,omitempty"`
	GDPPerCapitaUSD        *int     `json:"gdp_per_capita_usd,omitempty"`
	Population             *int64   `json:"population,omitempty"`
	TradeOpennessIndex     *float64 `json:"trade_openness_index,omitempty"`
	EaseOfDoingBusinessRank *int    `json:"ease_of_doing_business_rank,omitempty"`
	DistanceBrazilKm       int      `json:"distance_brazil_km"`
	Latitude               *float64 `json:"latitude,omitempty"`
	Longitude              *float64 `json:"longitude,omitempty"`
	FlagEmoji              string   `json:"flag_emoji,omitempty"`
	CurrencyCode           string   `json:"currency_code,omitempty"`
	Languages              []string `json:"languages,omitempty"`
}

// MarketData representa dados de mercado para um NCM × País
type MarketData struct {
	NCM              string  `json:"ncm"`
	CountryCode      string  `json:"country_code"`
	Year             int     `json:"year"`
	Month            int     `json:"month"`
	TotalValueUSD    float64 `json:"total_value_usd"`
	TotalWeightKg    float64 `json:"total_weight_kg"`
	AvgPricePerKgUSD float64 `json:"avg_price_per_kg_usd"`
	TransactionCount int     `json:"transaction_count"`
	GrowthRatePct    float64 `json:"growth_rate_pct"` // Comparado com período anterior
}

// ScoringWeights pesos para o algoritmo de scoring
type ScoringWeights struct {
	MarketSize  float64 // Peso do tamanho do mercado (default: 0.40)
	GrowthRate  float64 // Peso da taxa de crescimento (default: 0.30)
	PricePerKg  float64 // Peso do preço por kg (default: 0.20)
	Distance    float64 // Peso da distância (default: 0.10)
}

// DefaultScoringWeights retorna os pesos padrão do algoritmo
func DefaultScoringWeights() ScoringWeights {
	return ScoringWeights{
		MarketSize:  0.40,
		GrowthRate:  0.30,
		PricePerKg:  0.20,
		Distance:    0.10,
	}
}

// ValidateSimulatorRequest valida a requisição do simulador
func (r *SimulatorRequest) ValidateSimulatorRequest() error {
	// NCM deve ter exatamente 8 dígitos
	if len(r.NCM) != 8 {
		return ErrInvalidNCM
	}

	// Validar que NCM contém apenas dígitos
	for _, c := range r.NCM {
		if c < '0' || c > '9' {
			return ErrInvalidNCM
		}
	}

	// Volume deve ser positivo (se fornecido)
	if r.VolumeKg != nil && *r.VolumeKg <= 0 {
		return ErrInvalidVolume
	}

	// MaxResults deve estar entre 1 e 50
	if r.MaxResults < 0 || r.MaxResults > 50 {
		return ErrInvalidMaxResults
	}

	// Se MaxResults não fornecido, usar default
	if r.MaxResults == 0 {
		r.MaxResults = 10
	}

	return nil
}

// CalculateScore calcula o score geral de uma recomendação
func (d *DestinationRecommendation) CalculateScore(weights ScoringWeights, maxMarketSize, maxGrowthRate, maxPrice, maxDistance float64) float64 {
	// Normaliza cada métrica (0-1)
	marketSizeScore := 0.0
	if maxMarketSize > 0 {
		marketSizeScore = float64(d.MarketSizeUSD) / maxMarketSize
	}

	growthRateScore := 0.0
	if maxGrowthRate > 0 {
		growthRateScore = d.GrowthRatePct / maxGrowthRate
	}

	priceScore := 0.0
	if maxPrice > 0 {
		priceScore = d.PricePerKgUSD / maxPrice
	}

	// Distância: menor é melhor (inverter)
	distanceScore := 0.0
	if maxDistance > 0 {
		distanceScore = 1.0 - (float64(d.DistanceKm) / maxDistance)
	}

	// Score ponderado (0-1)
	score := (marketSizeScore * weights.MarketSize) +
		(growthRateScore * weights.GrowthRate) +
		(priceScore * weights.PricePerKg) +
		(distanceScore * weights.Distance)

	// Converte para escala 0-10
	return score * 10
}

// GetDemandLevel retorna o nível de demanda baseado no tamanho do mercado
func (d *DestinationRecommendation) GetDemandLevel() string {
	if d.MarketSizeUSD > 100_000_000 { // > 100M USD
		return "Alto"
	}
	if d.MarketSizeUSD > 10_000_000 { // > 10M USD
		return "Médio"
	}
	return "Baixo"
}

// GetRecommendationReason gera explicação do score
func (d *DestinationRecommendation) GetRecommendationReason(weights ScoringWeights) string {
	if d.Score >= 8.0 {
		return "Mercado altamente atrativo com grande potencial de crescimento e demanda consolidada"
	}
	if d.Score >= 6.0 {
		return "Mercado promissor com bom equilíbrio entre demanda, crescimento e custos logísticos"
	}
	if d.Score >= 4.0 {
		return "Mercado em desenvolvimento com oportunidades emergentes"
	}
	return "Mercado em fase inicial ou com barreiras significativas"
}
