-- Migration 0010: Simulator Tables
-- Creates tables for destination recommendation simulator

-- ============================================================================
-- COUNTRIES METADATA TABLE
-- ============================================================================
-- Stores metadata about countries for destination recommendations
CREATE TABLE IF NOT EXISTS public.countries_metadata (
  code TEXT PRIMARY KEY,                    -- ISO 3166-1 alpha-2 code (BR, US, CN)
  name_pt TEXT NOT NULL,                    -- Nome em português
  name_en TEXT NOT NULL,                    -- Name in English
  region TEXT NOT NULL,                     -- Region (Americas, Europe, Asia, Africa, Oceania)
  subregion TEXT,                          -- Subregion (South America, Western Europe, etc)

  -- Economic indicators
  gdp_usd BIGINT,                          -- GDP in USD (latest year)
  gdp_per_capita_usd INTEGER,              -- GDP per capita in USD
  population BIGINT,                        -- Population (latest year)

  -- Trade indicators
  trade_openness_index DECIMAL(5,2),       -- Trade openness (0-100)
  ease_of_doing_business_rank INTEGER,     -- World Bank ranking

  -- Geographic data
  distance_brazil_km INTEGER NOT NULL,     -- Distance from Brazil capital (km)
  latitude DECIMAL(10,7),                  -- Latitude
  longitude DECIMAL(10,7),                 -- Longitude

  -- Metadata
  flag_emoji TEXT,                         -- Flag emoji (🇧🇷)
  currency_code TEXT,                      -- Currency code (USD, EUR, BRL)
  languages TEXT[],                        -- Languages spoken

  -- Audit fields
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

  CONSTRAINT valid_distance CHECK (distance_brazil_km >= 0),
  CONSTRAINT valid_gdp CHECK (gdp_usd IS NULL OR gdp_usd >= 0),
  CONSTRAINT valid_population CHECK (population IS NULL OR population >= 0)
);

-- Indexes for countries_metadata
CREATE INDEX IF NOT EXISTS idx_countries_region
  ON public.countries_metadata(region);

CREATE INDEX IF NOT EXISTS idx_countries_distance
  ON public.countries_metadata(distance_brazil_km);

CREATE INDEX IF NOT EXISTS idx_countries_gdp
  ON public.countries_metadata(gdp_usd DESC NULLS LAST);

COMMENT ON TABLE public.countries_metadata IS 'Metadata about countries for destination recommendations';
COMMENT ON COLUMN public.countries_metadata.code IS 'ISO 3166-1 alpha-2 country code';
COMMENT ON COLUMN public.countries_metadata.distance_brazil_km IS 'Distance from Brasília to country capital';
COMMENT ON COLUMN public.countries_metadata.trade_openness_index IS 'Trade openness index (0-100, higher = more open)';

-- ============================================================================
-- COMEXSTAT CACHE TABLE
-- ============================================================================
-- Backup cache for ComexStat API data (fallback when API unavailable)
CREATE TABLE IF NOT EXISTS public.comexstat_cache (
  id BIGSERIAL PRIMARY KEY,

  -- Cache key components
  type TEXT NOT NULL,                      -- 'export' or 'import'
  year INTEGER NOT NULL,                   -- Year (2020-2025)
  month INTEGER NOT NULL,                  -- Month (1-12)
  ncm TEXT,                                -- NCM code (8 digits, NULL = all)
  country_code TEXT,                       -- Country ISO code (NULL = all)

  -- Aggregated data
  total_value_usd DECIMAL(18,2),          -- Total trade value in USD
  total_weight_kg DECIMAL(18,2),          -- Total weight in kg
  avg_price_per_kg_usd DECIMAL(10,4),     -- Average price per kg
  transaction_count INTEGER,               -- Number of transactions

  -- Raw data (compressed)
  raw_data JSONB,                         -- Full API response (compressed)

  -- Cache metadata
  api_status_code INTEGER,                 -- HTTP status from API
  api_response_time_ms INTEGER,            -- API response time in ms

  -- Audit fields
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  expires_at TIMESTAMPTZ NOT NULL,         -- TTL (7 days historic, 6h current month)
  hit_count INTEGER NOT NULL DEFAULT 0,    -- Number of cache hits

  CONSTRAINT valid_year CHECK (year >= 2020 AND year <= 2030),
  CONSTRAINT valid_month CHECK (month >= 1 AND month <= 12),
  CONSTRAINT valid_ncm CHECK (ncm IS NULL OR LENGTH(ncm) = 8),
  CONSTRAINT valid_type CHECK (type IN ('export', 'import'))
);

-- Unique constraint for cache key
CREATE UNIQUE INDEX IF NOT EXISTS idx_comexstat_cache_key
  ON public.comexstat_cache(type, year, month, COALESCE(ncm, ''), COALESCE(country_code, ''));

-- Index for expiration cleanup
CREATE INDEX IF NOT EXISTS idx_comexstat_cache_expires
  ON public.comexstat_cache(expires_at);

-- Index for frequent queries
CREATE INDEX IF NOT EXISTS idx_comexstat_cache_ncm
  ON public.comexstat_cache(ncm)
  WHERE ncm IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_comexstat_cache_country
  ON public.comexstat_cache(country_code)
  WHERE country_code IS NOT NULL;

-- Index for hit tracking
CREATE INDEX IF NOT EXISTS idx_comexstat_cache_hits
  ON public.comexstat_cache(hit_count DESC);

-- GIN index for JSONB queries
CREATE INDEX IF NOT EXISTS idx_comexstat_cache_raw_data
  ON public.comexstat_cache USING gin(raw_data);

COMMENT ON TABLE public.comexstat_cache IS 'L3 cache for ComexStat API data (fallback + historical backup)';
COMMENT ON COLUMN public.comexstat_cache.expires_at IS 'TTL: 7 days for historical data, 6h for current month';
COMMENT ON COLUMN public.comexstat_cache.hit_count IS 'Tracks cache popularity for eviction strategy';
COMMENT ON COLUMN public.comexstat_cache.raw_data IS 'Full API response in JSONB format (use for complex queries)';

-- ============================================================================
-- SIMULATOR RECOMMENDATIONS TABLE
-- ============================================================================
-- Stores simulator recommendations for analytics
CREATE TABLE IF NOT EXISTS public.simulator_recommendations (
  id BIGSERIAL PRIMARY KEY,

  -- Request data
  ncm TEXT NOT NULL,                       -- NCM code requested
  volume_kg DECIMAL(18,2),                 -- Volume requested (optional)
  user_id TEXT,                            -- User ID (if authenticated)
  session_id TEXT,                         -- Session ID (anonymous users)
  ip_address INET,                         -- IP address for rate limiting

  -- Response data
  recommendations JSONB NOT NULL,          -- Array of destination recommendations
  total_destinations INTEGER NOT NULL,     -- Number of destinations returned
  processing_time_ms INTEGER,              -- Processing time in ms

  -- Cache metadata
  cache_hit BOOLEAN NOT NULL DEFAULT false, -- Was this a cache hit?
  cache_level TEXT,                        -- 'l1', 'l2', 'l3', 'external'

  -- Audit fields
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

  CONSTRAINT valid_ncm_length CHECK (LENGTH(ncm) = 8),
  CONSTRAINT valid_volume CHECK (volume_kg IS NULL OR volume_kg > 0),
  CONSTRAINT valid_total_destinations CHECK (total_destinations >= 0)
);

-- Indexes for simulator_recommendations
CREATE INDEX IF NOT EXISTS idx_simulator_recommendations_ncm
  ON public.simulator_recommendations(ncm);

CREATE INDEX IF NOT EXISTS idx_simulator_recommendations_user
  ON public.simulator_recommendations(user_id)
  WHERE user_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_simulator_recommendations_session
  ON public.simulator_recommendations(session_id)
  WHERE session_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_simulator_recommendations_ip
  ON public.simulator_recommendations(ip_address, created_at);

CREATE INDEX IF NOT EXISTS idx_simulator_recommendations_created
  ON public.simulator_recommendations(created_at DESC);

-- GIN index for JSONB recommendations
CREATE INDEX IF NOT EXISTS idx_simulator_recommendations_data
  ON public.simulator_recommendations USING gin(recommendations);

COMMENT ON TABLE public.simulator_recommendations IS 'Analytics tracking for simulator usage';
COMMENT ON COLUMN public.simulator_recommendations.cache_level IS 'Which cache level served this request (for performance analysis)';

-- ============================================================================
-- FUNCTIONS
-- ============================================================================

-- Function to increment cache hit count
CREATE OR REPLACE FUNCTION increment_comexstat_cache_hit(
  p_type TEXT,
  p_year INTEGER,
  p_month INTEGER,
  p_ncm TEXT DEFAULT NULL,
  p_country_code TEXT DEFAULT NULL
)
RETURNS VOID AS $$
BEGIN
  UPDATE public.comexstat_cache
  SET
    hit_count = hit_count + 1,
    updated_at = now()
  WHERE
    type = p_type
    AND year = p_year
    AND month = p_month
    AND COALESCE(ncm, '') = COALESCE(p_ncm, '')
    AND COALESCE(country_code, '') = COALESCE(p_country_code, '');
END;
$$ LANGUAGE plpgsql;

-- Function to cleanup expired cache
CREATE OR REPLACE FUNCTION cleanup_expired_comexstat_cache()
RETURNS INTEGER AS $$
DECLARE
  deleted_count INTEGER;
BEGIN
  DELETE FROM public.comexstat_cache
  WHERE expires_at < now();

  GET DIAGNOSTICS deleted_count = ROW_COUNT;
  RETURN deleted_count;
END;
$$ LANGUAGE plpgsql;

-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = now();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Triggers for updated_at
CREATE TRIGGER update_countries_metadata_updated_at
  BEFORE UPDATE ON public.countries_metadata
  FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_comexstat_cache_updated_at
  BEFORE UPDATE ON public.comexstat_cache
  FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ============================================================================
-- INITIAL DATA SEEDING
-- ============================================================================
-- Seed top 50 countries (will be replaced by proper data loader)

-- Top 10 export destinations for Brazil (placeholder data)
INSERT INTO public.countries_metadata (code, name_pt, name_en, region, subregion, distance_brazil_km, gdp_usd, population) VALUES
  ('CN', 'China', 'China', 'Asia', 'Eastern Asia', 17500, 17963170000000, 1439323776),
  ('US', 'Estados Unidos', 'United States', 'Americas', 'Northern America', 7500, 25462700000000, 331002651),
  ('AR', 'Argentina', 'Argentina', 'Americas', 'South America', 2000, 487227000000, 45195774),
  ('NL', 'Países Baixos', 'Netherlands', 'Europe', 'Western Europe', 9500, 1012847000000, 17134872),
  ('CL', 'Chile', 'Chile', 'Americas', 'South America', 3500, 317059000000, 19116201),
  ('DE', 'Alemanha', 'Germany', 'Europe', 'Western Europe', 9800, 4259935000000, 83783942),
  ('JP', 'Japão', 'Japan', 'Asia', 'Eastern Asia', 18500, 4937420000000, 126476461),
  ('IN', 'Índia', 'India', 'Asia', 'Southern Asia', 14000, 3173398000000, 1380004385),
  ('MX', 'México', 'Mexico', 'Americas', 'Central America', 6500, 1293130000000, 128932753),
  ('ES', 'Espanha', 'Spain', 'Europe', 'Southern Europe', 8500, 1425865000000, 46754778)
ON CONFLICT (code) DO NOTHING;

COMMENT ON TABLE public.countries_metadata IS 'Metadata about top 50 trading partner countries for destination recommendations';
COMMENT ON TABLE public.comexstat_cache IS 'L3 cache: ComexStat API data backup for fallback and historical analysis';
COMMENT ON TABLE public.simulator_recommendations IS 'Analytics: Tracks all simulator requests for usage analysis and ML training';
