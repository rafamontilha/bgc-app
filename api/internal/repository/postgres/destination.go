package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"bgc-app/internal/business/destination"

	"github.com/lib/pq"
)

// DestinationRepository implementa destination.Repository para PostgreSQL
type DestinationRepository struct {
	db *sql.DB
}

// NewDestinationRepository cria uma nova instância do repository
func NewDestinationRepository(db *sql.DB) *DestinationRepository {
	return &DestinationRepository{
		db: db,
	}
}

// GetCountryMetadata busca metadados de um país específico
func (r *DestinationRepository) GetCountryMetadata(ctx context.Context, countryCode string) (*destination.CountryMetadata, error) {
	query := `
		SELECT
			code, name_pt, name_en, region, subregion,
			gdp_usd, gdp_per_capita_usd, population,
			trade_openness_index, ease_of_doing_business_rank,
			distance_brazil_km, latitude, longitude,
			flag_emoji, currency_code, languages
		FROM public.countries_metadata
		WHERE code = $1
	`

	var country destination.CountryMetadata
	var languages pq.StringArray

	err := r.db.QueryRowContext(ctx, query, countryCode).Scan(
		&country.Code,
		&country.NamePt,
		&country.NameEn,
		&country.Region,
		&country.Subregion,
		&country.GDPUSD,
		&country.GDPPerCapitaUSD,
		&country.Population,
		&country.TradeOpennessIndex,
		&country.EaseOfDoingBusinessRank,
		&country.DistanceBrazilKm,
		&country.Latitude,
		&country.Longitude,
		&country.FlagEmoji,
		&country.CurrencyCode,
		&languages,
	)

	if err == sql.ErrNoRows {
		return nil, destination.ErrCountryNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get country metadata: %w", err)
	}

	country.Languages = []string(languages)

	return &country, nil
}

// GetAllCountries busca todos os países disponíveis
func (r *DestinationRepository) GetAllCountries(ctx context.Context) ([]destination.CountryMetadata, error) {
	query := `
		SELECT
			code, name_pt, name_en, region, subregion,
			gdp_usd, gdp_per_capita_usd, population,
			trade_openness_index, ease_of_doing_business_rank,
			distance_brazil_km, latitude, longitude,
			flag_emoji, currency_code, languages
		FROM public.countries_metadata
		ORDER BY name_pt
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query countries: %w", err)
	}
	defer rows.Close()

	var countries []destination.CountryMetadata

	for rows.Next() {
		var country destination.CountryMetadata
		var languages pq.StringArray

		err := rows.Scan(
			&country.Code,
			&country.NamePt,
			&country.NameEn,
			&country.Region,
			&country.Subregion,
			&country.GDPUSD,
			&country.GDPPerCapitaUSD,
			&country.Population,
			&country.TradeOpennessIndex,
			&country.EaseOfDoingBusinessRank,
			&country.DistanceBrazilKm,
			&country.Latitude,
			&country.Longitude,
			&country.FlagEmoji,
			&country.CurrencyCode,
			&languages,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan country: %w", err)
		}

		country.Languages = []string(languages)
		countries = append(countries, country)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating countries: %w", err)
	}

	return countries, nil
}

// GetMarketDataByNCM busca dados de mercado por NCM
// Retorna dados agregados dos últimos 12 meses por país
func (r *DestinationRepository) GetMarketDataByNCM(ctx context.Context, ncm string, year, month int) ([]destination.MarketData, error) {
	// Busca dados dos últimos 12 meses
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC).AddDate(0, -12, 0)

	query := `
		WITH market_data AS (
			SELECT
				$1 as ncm,
				co_pais as country_code,
				co_ano as year,
				co_mes as month,
				SUM(vl_fob) as total_value_usd,
				SUM(kg_liquido) as total_weight_kg,
				COUNT(*) as transaction_count
			FROM stg.exportacao
			WHERE
				SUBSTRING(co_ncm, 1, 8) = $1
				AND co_ano >= $2
				AND (co_ano > $2 OR co_mes >= $3)
			GROUP BY co_pais, co_ano, co_mes
		),
		aggregated AS (
			SELECT
				ncm,
				country_code,
				MAX(year) as year,
				MAX(month) as month,
				SUM(total_value_usd) as total_value_usd,
				SUM(total_weight_kg) as total_weight_kg,
				SUM(transaction_count) as transaction_count,
				CASE
					WHEN SUM(total_weight_kg) > 0
					THEN SUM(total_value_usd) / SUM(total_weight_kg)
					ELSE 0
				END as avg_price_per_kg_usd
			FROM market_data
			GROUP BY ncm, country_code
		),
		previous_period AS (
			SELECT
				country_code,
				SUM(total_value_usd) as prev_value
			FROM market_data
			WHERE
				year = $2 - 1
				OR (year = $2 AND month < $3)
			GROUP BY country_code
		)
		SELECT
			a.ncm,
			a.country_code,
			a.year,
			a.month,
			a.total_value_usd,
			a.total_weight_kg,
			a.avg_price_per_kg_usd,
			a.transaction_count,
			CASE
				WHEN p.prev_value > 0
				THEN ((a.total_value_usd - p.prev_value) / p.prev_value * 100)
				ELSE 0
			END as growth_rate_pct
		FROM aggregated a
		LEFT JOIN previous_period p ON a.country_code = p.country_code
		WHERE a.total_value_usd > 0
		ORDER BY a.total_value_usd DESC
		LIMIT 100
	`

	rows, err := r.db.QueryContext(ctx, query, ncm, startDate.Year(), int(startDate.Month()))
	if err != nil {
		return nil, fmt.Errorf("failed to query market data: %w", err)
	}
	defer rows.Close()

	var marketData []destination.MarketData

	for rows.Next() {
		var data destination.MarketData

		err := rows.Scan(
			&data.NCM,
			&data.CountryCode,
			&data.Year,
			&data.Month,
			&data.TotalValueUSD,
			&data.TotalWeightKg,
			&data.AvgPricePerKgUSD,
			&data.TransactionCount,
			&data.GrowthRatePct,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan market data: %w", err)
		}

		marketData = append(marketData, data)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating market data: %w", err)
	}

	if len(marketData) == 0 {
		return nil, destination.ErrNoDataAvailable
	}

	return marketData, nil
}

// GetMarketDataByNCMAndCountry busca dados de mercado por NCM e país específico
func (r *DestinationRepository) GetMarketDataByNCMAndCountry(ctx context.Context, ncm, countryCode string, year, month int) (*destination.MarketData, error) {
	query := `
		SELECT
			$1 as ncm,
			co_pais as country_code,
			$2 as year,
			$3 as month,
			SUM(vl_fob) as total_value_usd,
			SUM(kg_liquido) as total_weight_kg,
			CASE
				WHEN SUM(kg_liquido) > 0
				THEN SUM(vl_fob) / SUM(kg_liquido)
				ELSE 0
			END as avg_price_per_kg_usd,
			COUNT(*) as transaction_count,
			0 as growth_rate_pct
		FROM stg.exportacao
		WHERE
			SUBSTRING(co_ncm, 1, 8) = $1
			AND co_pais = $4
			AND co_ano = $2
			AND co_mes = $3
		GROUP BY co_pais
	`

	var data destination.MarketData

	err := r.db.QueryRowContext(ctx, query, ncm, year, month, countryCode).Scan(
		&data.NCM,
		&data.CountryCode,
		&data.Year,
		&data.Month,
		&data.TotalValueUSD,
		&data.TotalWeightKg,
		&data.AvgPricePerKgUSD,
		&data.TransactionCount,
		&data.GrowthRatePct,
	)

	if err == sql.ErrNoRows {
		return nil, destination.ErrNoDataAvailable
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get market data: %w", err)
	}

	return &data, nil
}

// SaveRecommendation salva uma recomendação para analytics
func (r *DestinationRepository) SaveRecommendation(ctx context.Context, req destination.SimulatorRequest, resp destination.SimulatorResponse, userID, sessionID string, ipAddress string, cacheHit bool, cacheLevel string) error {
	query := `
		INSERT INTO public.simulator_recommendations (
			ncm, volume_kg, user_id, session_id, ip_address,
			recommendations, total_destinations, processing_time_ms,
			cache_hit, cache_level, created_at
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8,
			$9, $10, now()
		)
	`

	// Converte recommendations para JSONB
	recommendationsJSON, err := json.Marshal(resp.Destinations)
	if err != nil {
		return fmt.Errorf("failed to marshal recommendations: %w", err)
	}

	var userIDPtr, sessionIDPtr, ipAddressPtr, cacheLevelPtr *string
	if userID != "" {
		userIDPtr = &userID
	}
	if sessionID != "" {
		sessionIDPtr = &sessionID
	}
	if ipAddress != "" {
		ipAddressPtr = &ipAddress
	}
	if cacheLevel != "" {
		cacheLevelPtr = &cacheLevel
	}

	_, err = r.db.ExecContext(ctx, query,
		req.NCM,
		req.VolumeKg,
		userIDPtr,
		sessionIDPtr,
		ipAddressPtr,
		recommendationsJSON,
		resp.Metadata.TotalDestinations,
		resp.Metadata.ProcessingTimeMs,
		cacheHit,
		cacheLevelPtr,
	)

	if err != nil {
		return fmt.Errorf("failed to save recommendation: %w", err)
	}

	return nil
}
